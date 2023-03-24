package security

type ManagerCertificates interface {
	VerifyCertificate() bool
	GetCertificateCA() error
	GetCertificateHost() error
}
