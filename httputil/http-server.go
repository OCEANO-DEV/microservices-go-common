package httputil

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/oceano-dev/microservices-go-common/services"
)

type HttpServer interface {
	RunTLSServer() (*http.Server, error)
}

type httpServer struct {
	config  *config.Config
	router  *gin.Engine
	service services.CertificatesService
}

var srv *http.Server

func NewHttpServer(
	config *config.Config,
	router *gin.Engine,
	service services.CertificatesService,
) *httpServer {
	return &httpServer{
		config:  config,
		router:  router,
		service: service,
	}
}

func (s *httpServer) RunTLSServer() (*http.Server, error) {
	var err error
	if srv == nil {
		srv = s.mountTLSServer()

		go func() {
			if err = srv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				log.Fatalf("err: %s\n", err)
			}
		}()
	}

	return srv, err
}

func (s *httpServer) mountTLSServer() *http.Server {
	return &http.Server{
		Addr:         s.config.ListenPort,
		Handler:      s.router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			GetCertificate: s.service.GetLocalCertificate,
			RootCAs:        s.service.GetLocalCertificateCA(),
		},
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
}

// func (s *httpServer) getLocalCertificate(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
// 	certPath, keyPath := s.service.GetPathsCertificateAndKey()
// 	if !helpers.FileExists(certPath) || !helpers.FileExists(keyPath) {
// 		return nil, errors.New("certificate not found")
// 	}

// 	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}

// 	return &cert, nil
// }
