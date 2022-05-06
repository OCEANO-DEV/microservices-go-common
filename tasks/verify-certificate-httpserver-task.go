package tasks

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/oceano-dev/microservices-go-common/httputil"
	"github.com/oceano-dev/microservices-go-common/security"
	"github.com/oceano-dev/microservices-go-common/services"
	trace "github.com/oceano-dev/microservices-go-common/trace/otel"
)

type VerifyCertificateWithHttpServerTask struct {
	config              *config.Config
	managerCertificates *security.ManagerCertificates
	emailService        services.EmailService
	httputil            httputil.HttpServer
}

func NewVerifyCertificateWithHttpServerTask(
	config *config.Config,
	managerCertificates *security.ManagerCertificates,
	emailService services.EmailService,
	httputil httputil.HttpServer,
) *VerifyCertificateWithHttpServerTask {
	return &VerifyCertificateWithHttpServerTask{
		config:              config,
		managerCertificates: managerCertificates,
		emailService:        emailService,
		httputil:            httputil,
	}
}

var srv *http.Server

func (task *VerifyCertificateWithHttpServerTask) ReloadCertificate(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				_, span := trace.NewSpan(ctx, "VerifyCertificateTask.ReloadCertificate")
				defer span.End()

				certIsValid := task.managerCertificates.VerifiyLocalCertificateIsValid()
				if !certIsValid {
					err := task.managerCertificates.GetCertificate()
					if err != nil {
						msg := fmt.Sprintln("EmailService - certificate error: ", err)
						err := task.emailService.SendSupportMessage(msg)
						if err != nil {
							log.Println(err)
						}
						log.Println(msg)
						ticker.Reset(60 * time.Second)
						break
					}
				}
				fmt.Printf("certificate verified successfully: %s\n", time.Now().UTC())

				if srv == nil {
					if certIsValid {
						srvNew, err := task.httputil.RunTLSServer()
						if err != nil {
							log.Fatal("http server error: ", err)
							break
						}

						srv = srvNew
						log.Printf("Listening on port %s", task.config.ListenPort)
					}
				}

				ticker.Reset(1 * time.Minute)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
