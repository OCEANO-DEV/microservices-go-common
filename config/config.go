package config

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/twinj/uuid"
)

type Config struct {
	Production      bool                  `json:"production"`
	AppName         string                `json:"appName"`
	ApiVersion      string                `json:"apiVersion"`
	Endpoint        string                `json:"endpoint"`
	ListenPort      string                `json:"listenPort"`
	Folders         []string              `json:"folders"`
	MongoDB         MongoDBConfig         `json:"mongodb"`
	Certificates    CertificatesConfig    `json:"certificates"`
	Token           TokenConfig           `json:"token"`
	SecurityKeys    SecurityKeysConfig    `json:"securityKeys"`
	SecurityRSAKeys SecurityRSAKeysConfig `json:"securityRSAKeys"`
	SMTPServer      SMTPConfig            `json:"smtpServer"`
	Company         CompanyConfig         `json:"company"`
	Prometheus      PrometheusConfig      `json:"prometheus"`
	MongoDbExporter MongoDbExporterConfig `json:"mongoDbExporter"`
	Nats            NatsConfig            `json:"nats"`
	Jaeger          JaegerConfig          `json:"jaeger"`
	GrpcServer      GrpcServerConfig      `json:"grpcServer"`
	EmailService    EmailServiceConfig    `json:"emailService"`
	Postgres        PostgresConfig        `json:"postgres"`
	Redis           RedisConfig           `json:"redis"`
}

type MongoDBConfig struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	Database    string `json:"database"`
	User        string `json:"user"`
	Password    string `json:"password"`
	MaxPoolSize int    `json:"maxPoolSize"`
}

type CertificatesConfig struct {
	FileName                   string `json:"filename"`
	HashPermissionEndPoint     string `json:"hashPermissionEndPoint"`
	PasswordPermissionEndPoint string `json:"passwordPermissionEndPoint"`
	EndPointGetCertificate     string `json:"endPointGetCertificate"`
	EndPointGetCertificateKey  string `json:"endPointGetCertificateKey"`
	MinutesToReloadCertificate int    `json:"minutesToReloadCertificate"`
}

type TokenConfig struct {
	Issuer                    string `json:"issuer"`
	MinutesToExpireToken      int    `json:"minutesToExpireToken"`
	HoursToExpireRefreshToken int    `json:"hoursToExpireRefreshToken"`
}

type SecurityKeysConfig struct {
	DaysToExpireKeys            int    `json:"daysToExpireKeys"`
	MinutesToRefreshPrivateKeys int    `json:"minutesToRefreshPrivateKeys"`
	MinutesToRefreshPublicKeys  int    `json:"minutesToRefreshPublicKeys"`
	SavePublicKeyToFile         bool   `json:"savePublicKeyToFile"`
	FileECPPublicKey            string `json:"fileECPPublicKey"`
	EndPointGetPublicKeys       string `json:"endPointGetPublicKeys"`
}

type SecurityRSAKeysConfig struct {
	DaysToExpireRSAKeys            int    `json:"daysToExpireRSAKeys"`
	MinutesToRefreshRSAPrivateKeys int    `json:"minutesToRefreshRSAPrivateKeys"`
	MinutesToRefreshRSAPublicKeys  int    `json:"minutesToRefreshRSAPublicKeys"`
	EndPointGetRSAPublicKeys       string `json:"endPointGetRSAPublicKeys"`
}

type CompanyConfig struct {
	Name              string `json:"name"`
	Address           string `json:"address"`
	AddressNumber     string `json:"addressNumber"`
	AddressComplement string `json:"addressComplement"`
	Locality          string `json:"locality"`
	Country           string `json:"country"`
	PostalCode        string `json:"postalCode"`
	Phone             string `json:"phone"`
	Email             string `json:"email"`
}

type SMTPConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	TLS          bool   `json:"tls"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	SupportEmail string `json:"supportEmail"`
}

type PrometheusConfig struct {
	PROMETHEUS_PUSHGATEWAY string `json:"prometheus_pushgateway"`
}

type MongoDbExporterConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type NatsConfig struct {
	Url         string `json:"url"`
	ClusterId   string `json:"clusterId"`
	ClientId    string `json:"clientId"`
	ConnectWait int    `json:"connectWait"`
	PubAckWait  int    `json:"pubAckWait"`
	Interval    int    `json:"interval"`
	MaxOut      int    `json:"maxOut"`
}

type JaegerConfig struct {
	JaegerEndpoint string `json:"jaegerEndpoint"`
	ServiceName    string `json:"serviceName"`
	ServiceVersion string `json:"serviceVersion"`
}

type GrpcServerConfig struct {
	Port              string `json:"port"`
	MaxConnectionIdle int    `json:"maxConnectionIdle"`
	MaxConnectionAge  int    `json:"maxConnectionAge"`
	Timeout           int    `json:"timeout"`
}

type EmailServiceConfig struct {
	Host string `json:"host"`
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	SSLMode  string `json:"sslMode"`
}

type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	Db       int    `json:"db"`
	PoolSize int    `json:"poolSize"`
}

func LoadConfig(production bool, path string) *Config {
	viper.AddConfigPath(path)
	viper.SetConfigName("config-dev")
	if production {
		viper.SetConfigName("config-prod")
	}
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	config := &Config{}
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshal config: %s", err))
	}

	config.Production = production

	if config.Production {
		config.Nats.ClientId += "_" + uuid.NewV4().String()

		fmt.Printf("ENVIRONMENT: production\n")
	}

	fmt.Printf("MONGO_HOST: %s\nMONGO_PORT: %s\n", config.MongoDB.Host, config.MongoDB.Port)
	fmt.Printf("NATS_URL: %s\n", config.Nats.Url)

	return config
}
