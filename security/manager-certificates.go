package security

import (
	"errors"
	"time"

	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/oceano-dev/microservices-go-common/helpers"
	"github.com/oceano-dev/microservices-go-common/services"

	"github.com/eapache/go-resiliency/breaker"
)

type managerCertificates struct {
	config  *config.Config
	service services.CertificatesService
}

var (
	certPath string
	keyPath  string
)

func NewManagerCertificates(
	config *config.Config,
	service services.CertificatesService,
) *managerCertificates {
	certPath, keyPath = service.GetPathsCertificateAndKey()
	return &managerCertificates{
		config:  config,
		service: service,
	}
}

func (m *managerCertificates) VerifyCertificate() bool {
	if helpers.FileExists(certPath) && helpers.FileExists(keyPath) {
		cert, err := m.service.ReadCertificate(certPath)
		if err != nil {
			return false
		}

		if cert == nil || cert.NotAfter.AddDate(0, 0, -7).Before(time.Now().UTC()) {
			return false
		}

		return true
	}

	return false
}

func (m *managerCertificates) GetCertificate() error {
	err := m.refreshCertificate()
	if err != nil {
		return err
	}

	return nil
}

func (m *managerCertificates) refreshCertificate() error {
	err := m.requestCertificate()
	if err != nil {
		return err
	}

	err = m.requestCertificateKey()
	if err != nil {
		return err
	}

	return nil
}

func (m managerCertificates) requestCertificate() error {
	b := breaker.New(3, 1, 5*time.Second)
	for {
		var cert []byte
		var err error
		err = b.Run(func() error {
			cert, err = m.service.GetCertificate()
			if err != nil {
				return err
			}

			return nil
		})

		switch err {
		case nil:
			if cert == nil {
				return errors.New("certificate not found")
			}

			err := helpers.CreateFile(cert, certPath)
			if err != nil {
				return err
			}

			return nil
		case breaker.ErrBreakerOpen:
			return err
		}
	}
}

func (m *managerCertificates) requestCertificateKey() error {
	b := breaker.New(3, 1, 5*time.Second)
	for {
		var key []byte
		var err error
		err = b.Run(func() error {
			key, err = m.service.GetCertificateKey()
			if err != nil {
				return err
			}

			return nil
		})

		switch err {
		case nil:
			if key == nil {
				return errors.New("certificate key not found")
			}

			err := helpers.CreateFile(key, keyPath)
			if err != nil {
				return err
			}

			return nil
		case breaker.ErrBreakerOpen:
			return err
		}
	}
}

func (m *managerCertificates) getCertificateKey() ([]byte, error) {
	key, err := m.service.GetCertificateKey()
	if err != nil {
		return nil, err
	}

	return key, nil
}
