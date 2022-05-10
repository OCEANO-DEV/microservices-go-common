package services

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/oceano-dev/microservices-go-common/config"
	"github.com/oceano-dev/microservices-go-common/helpers"
)

type CertificatesService interface {
	GetCertificate() ([]byte, error)
	GetCertificateKey() ([]byte, error)
	GetPathsCertificateAndKey() (string, string)
	ReadCertificate(pathCertificate string) (*x509.Certificate, error)
}

type certificatesService struct {
	config *config.Config
}

var (
	certPath string
	keyPath  string
)

func NewCertificatesService(
	config *config.Config,
) *certificatesService {
	certPath = fmt.Sprintf("certs/%s.crt", config.Certificates.FileName)
	keyPath = fmt.Sprintf("certs/%s.key", config.Certificates.FileName)
	return &certificatesService{
		config: config,
	}
}

func (s *certificatesService) GetCertificate() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data, err := s.requestCertificate(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *certificatesService) GetCertificateKey() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	data, err := s.requestCertificateKey(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *certificatesService) GetPathsCertificateAndKey() (string, string) {
	if !helpers.FileExists(certPath) || !helpers.FileExists(keyPath) {
		return "", ""
	}

	return certPath, keyPath
}

func (s *certificatesService) ReadCertificate(pathCertificate string) (*x509.Certificate, error) {
	data, err := ioutil.ReadFile(pathCertificate)
	if err != nil {
		os.Exit(1)
		return nil, fmt.Errorf("read Certificate file error")
	}

	pemBlock, _ := pem.Decode(data)
	if pemBlock == nil {
		return nil, fmt.Errorf("decode Certificate error")
	}

	cert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func (s *certificatesService) requestCertificate(ctx context.Context) ([]byte, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	hash := base64.StdEncoding.EncodeToString([]byte(s.config.Certificates.PasswordPermissionEndPoint))
	endPoint := fmt.Sprintf("%s/%s", s.config.Certificates.EndPointGetCertificate, hash)
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

func (s *certificatesService) requestCertificateKey(ctx context.Context) ([]byte, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	hash := base64.StdEncoding.EncodeToString([]byte(s.config.Certificates.PasswordPermissionEndPoint))
	endPoint := fmt.Sprintf("%s/%s", s.config.Certificates.EndPointGetCertificateKey, hash)
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
