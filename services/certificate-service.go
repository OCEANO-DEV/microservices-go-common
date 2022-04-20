package services

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/oceano-dev/microservices-go-common/config"
)

type CertificateService interface {
	GetCertificate() ([]byte, error)
	GetCertificateKey() ([]byte, error)
}

type certificateService struct {
	config *config.Config
}

func NewCertificateService(
	config *config.Config,
) *certificateService {
	return &certificateService{
		config: config,
	}
}

func (s *certificateService) GetCertificate() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data, err := s.requestCertificate(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *certificateService) GetCertificateKey() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data, err := s.requestCertificateKey(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *certificateService) requestCertificate(ctx context.Context) ([]byte, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	endPoint := fmt.Sprintf("%s/hash=%s", s.config.Certificates.EndPointGetCertificate, s.config.Certificates.HashPermissionEndPoint)
	request, err := http.NewRequestWithContext(ctx, "GET", endPoint, nil)
	if err != nil {
		log.Println("request:", err)
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		log.Println("response:", err)
		return nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("data parse:", err)
		return nil, err
	}

	return data, nil
}

func (s *certificateService) requestCertificateKey(ctx context.Context) ([]byte, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	endPoint := fmt.Sprintf("%s/hash=%s", s.config.Certificates.EndPointGetCertificateKey, s.config.Certificates.HashPermissionEndPoint)
	request, err := http.NewRequestWithContext(ctx, "GET", endPoint, nil)
	if err != nil {
		log.Println("request:", err)
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		log.Println("response:", err)
		return nil, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("data parse:", err)
		return nil, err
	}

	return data, nil
}
