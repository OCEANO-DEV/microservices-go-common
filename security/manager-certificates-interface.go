package security

type ManagerCertificates interface {
	VerifyCertificate() bool
	VerifyCertificateCA() bool
	GetCertificate() error
	GetCertificateCA() error
}
