package httputil

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/oceano-dev/microservices-go-common/security"
)

type HttpServer interface {
	RunTLSServer() error
}

type httpServer struct {
	config              *config.Config
	router              *gin.Engine
	managerCertificates *security.ManagerCertificates
}

var mux sync.Mutex
var srv *http.Server

func NewHttpServer(
	config *config.Config,
	router *gin.Engine,
	managerCertificates *security.ManagerCertificates,
) *httpServer {
	return &httpServer{
		config:              config,
		router:              router,
		managerCertificates: managerCertificates,
	}
}

func (s *httpServer) RunTLSServer() error {
	var err error
	if srv == nil {
		srv = s.mountTLSServer()

		go func() {
			if err = srv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				log.Fatalf("err: %s\n", err)
			}
		}()
	}

	return err
}

func (s *httpServer) mountTLSServer() *http.Server {
	return &http.Server{
		Addr:    s.config.ListenPort,
		Handler: s.router,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
			GetCertificate:           s.getCertificate,
		},
	}
}

func (s *httpServer) getCertificate(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	mux.Lock()
	defer mux.Unlock()

	pathCert, pathKey := s.managerCertificates.GetPathsCertificateAndKey()
	cert, err := tls.LoadX509KeyPair(pathCert, pathKey)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &cert, nil
}
