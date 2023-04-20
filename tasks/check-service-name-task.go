package tasks

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	consul "github.com/hashicorp/consul/api"
	"github.com/oceano-dev/microservices-go-common/config"
	trace "github.com/oceano-dev/microservices-go-common/trace/otel"

	parse "github.com/oceano-dev/microservices-go-common/consul"
)

type checkServiceNameTask struct{}

func NewCheckServiceNameTask() *checkServiceNameTask {
	return &checkServiceNameTask{}
}

func (task *checkServiceNameTask) ReloadServiceName(
	ctx context.Context,
	config *config.Config,
	consulClient *consul.Client,
	serviceName string,
	consulParse parse.ConsulParse,
	servicesNameDone chan bool) {
	ticker := time.NewTicker(2500 * time.Millisecond)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				_, span := trace.NewSpan(ctx, "checkServiceNameTask.ReloadServiceName")
				defer span.End()

				services, _, err := consulClient.Catalog().Service(serviceName, "", nil)
				// services, err := consulClient.Agent().Services()
				if err != nil {
					fmt.Printf("failed to refresh service name %s. error: %s", serviceName, err)
					ticker.Reset(5 * time.Second)
					break
				}

				ok := task.updateEndPoint(serviceName, config, services, consulParse)

				ticker.Reset(time.Duration(config.SecondsToReloadServicesName) * time.Second)
				if ok {
					fmt.Printf("refresh service name %s successfully: %s\n", serviceName, time.Now().UTC())
					servicesNameDone <- ok
				} else {
					fmt.Printf("service name %s not found. Refresh was not successfully: %s\n", serviceName, time.Now().UTC())
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (task *checkServiceNameTask) updateEndPoint(
	serviceName string,
	config *config.Config,
	services []*consul.CatalogService,
	// services map[string]*consul.AgentService,
	consulParse parse.ConsulParse,
) bool {

	if len(services) <= 0 {
		return false
	}

	qtd := len(services)
	service := services[rand.Intn(qtd)-1]

	// service := services[serviceName]
	// if service == nil {
	// 	return false
	// }

	host := fmt.Sprintf("https://%s:%s", service.Address, strconv.Itoa(service.ServicePort))

	switch consulParse {
	case parse.CertificatesAndSecurityKeys:
		config.Certificates.EndPointGetCertificateCA = fmt.Sprintf("%s/%s", host, config.Certificates.APIPathCertificateCA)
		config.Certificates.EndPointGetCertificateHost = fmt.Sprintf("%s/%s", host, config.Certificates.APIPathCertificateHost)
		config.Certificates.EndPointGetCertificateHostKey = fmt.Sprintf("%s/%s", host, config.Certificates.APIPathCertificateHostKey)
		config.SecurityKeys.EndPointGetPublicKeys = fmt.Sprintf("%s/%s", host, config.SecurityKeys.APIPathPublicKeys)
		return true

	case parse.SecurityRSAKeys:
		config.SecurityRSAKeys.EndPointGetRSAPublicKeys = fmt.Sprintf("%s/%s", host, config.SecurityRSAKeys.APIPathRSAPublicKeys)
		return true

	case parse.EmailService:
		config.EmailService.Host = host
		return true

	default:
		return false
	}
}
