package security

import "crypto/x509"

type ManagerCertificates interface {
	VerifyCertificate() bool
	GetCertificate() error
	GetPathsCertificateAndKey() (string, string)
	ReadCertificate(pathCertificate string) (*x509.Certificate, error)
}
