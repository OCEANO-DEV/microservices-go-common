package tasks

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/oceano-dev/microservices-go-common/httputil"
	"github.com/oceano-dev/microservices-go-common/security"
	"github.com/oceano-dev/microservices-go-common/services"
	trace "github.com/oceano-dev/microservices-go-common/trace/otel"
)

type SecureHttpServerTask struct {
	config              *config.Config
	managerCertificates security.ManagerCertificates
	emailService        services.EmailService
	httputil            httputil.HttpServer
}

func NewSecureHttpServerTask(
	config *config.Config,
	managerCertificates security.ManagerCertificates,
	emailService services.EmailService,
	httputil httputil.HttpServer,
) *SecureHttpServerTask {
	return &SecureHttpServerTask{
		config:              config,
		managerCertificates: managerCertificates,
		emailService:        emailService,
		httputil:            httputil,
	}
}

var muxHttpServer sync.Mutex
var srv *http.Server

func (task *SecureHttpServerTask) Start(ctx context.Context) {
	ticker := time.NewTicker(2500 * time.Millisecond)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				_, span := trace.NewSpan(ctx, "SecureHttpServerTask.Start")
				defer span.End()

				certIsValid := task.managerCertificates.VerifyCertificate()
				if !certIsValid {
					err := task.managerCertificates.GetCertificateCA()
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
				//fmt.Printf("start secure http server successfully: %s\n", time.Now().UTC())

				if srv == nil {
					if certIsValid {
						srvNew, err := task.httputil.RunTLSServer()
						if err != nil {
							log.Fatal("http server error: ", err)
							break
						}

						muxHttpServer.Lock()
						srv = srvNew
						muxHttpServer.Unlock()

						log.Printf("Listening on port %s", task.config.ListenPort)
					}
				}

				ticker.Reset(time.Duration(task.config.Certificates.MinutesToReloadCertificate) * time.Minute)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
