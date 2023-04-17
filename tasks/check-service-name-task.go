package tasks

import (
	"context"
	"fmt"
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

				services, error := consulClient.Agent().Services()
				if error != nil {
					fmt.Printf("failed to refresh service name %s. error: %s", serviceName, error)
					ticker.Reset(5 * time.Second)
					break
				}

				ok := task.updateEndPoint(serviceName, config, services, consulParse)

				fmt.Printf("start refresh service name successfully: %s\n", time.Now().UTC())
				ticker.Reset(time.Duration(config.SecondsToReloadServicesName) * time.Second)
				if ok {
					servicesNameDone <- ok
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
	services map[string]*consul.AgentService,
	consulParse parse.ConsulParse,
) bool {
	service := services[serviceName]

	switch consulParse {
	case parse.CertificatesAndSecurityKeys:
		endPoint := fmt.Sprintf("%s:%s/%s", service.Address, strconv.Itoa(service.Port), config.Certificates.APIPathCertificateCA)
		config.Certificates.EndPointGetCertificateCA = endPoint

		endPoint = fmt.Sprintf("%s:%s/%s", service.Address, strconv.Itoa(service.Port), config.Certificates.APIPathCertificateHost)
		config.Certificates.EndPointGetCertificateHost = endPoint

		endPoint = fmt.Sprintf("%s:%s/%s", service.Address, strconv.Itoa(service.Port), config.Certificates.APIPathCertificateHostKey)
		config.Certificates.EndPointGetCertificateHostKey = endPoint

		endPoint = fmt.Sprintf("%s:%s/%s", service.Address, strconv.Itoa(service.Port), config.SecurityKeys.APIPathPublicKeys)
		config.SecurityKeys.EndPointGetPublicKeys = endPoint
		return true

	case parse.SecurityRSAKeys:
		endPoint := fmt.Sprintf("%s:%s/%s", service.Address, strconv.Itoa(service.Port), config.SecurityRSAKeys.APIPathRSAPublicKeys)
		config.SecurityRSAKeys.EndPointGetRSAPublicKeys = endPoint
		return true

	case parse.EmailService:
		endPoint := fmt.Sprintf("%s:%s/%s", service.Address, strconv.Itoa(service.Port), config.EmailService.Host)
		config.EmailService.Host = endPoint
		return true

	default:
		return false
	}
}
