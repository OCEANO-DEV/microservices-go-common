package security

type ManagerCertificates interface {
	VerifyCertificate() bool
	GetCertificate() error
	GetCertificateCA() error
}
