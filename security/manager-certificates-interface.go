package security

type ManagerCertificates interface {
	VerifyCertificate() bool
	GetCertificateCA() error
	GetCertificate() error
}
