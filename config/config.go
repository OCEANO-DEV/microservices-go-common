package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"github.com/twinj/uuid"
)

type Config struct {
	Production      bool                  `json:"production"`
	AppName         string                `json:"appName"`
	ApiVersion      string                `json:"apiVersion"`
	Endpoint        string                `json:"endpoint"`
	ListenPort      string                `json:"listenPort"`
	MongoDB         MongoDBConfig         `json:"mongodb"`
	Certificates    CertificatesConfig    `json:"certificates"`
	Token           TokenConfig           `json:"token"`
	SecurityKeys    SecurityKeysConfig    `json:"securityKeys"`
	SMTPServer      SMTPConfig            `json:"smtpServer"`
	Company         CompanyConfig         `json:"company"`
	Prometheus      PrometheusConfig      `json:"prometheus"`
	MongoDbExporter MongoDbExporterConfig `json:"mongoDbExporter"`
	Nats            NatsConfig            `json:"nats"`
	Jaeger          JaegerConfig          `json:"jaeger"`
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
	FileName string `json:"filename"`
}

type TokenConfig struct {
	Issuer                    string `json:"issuer"`
	MinutesToExpireToken      int    `json:"minutesToExpireToken"`
	HoursToExpireRefreshToken int    `json:"hoursToExpireRefreshToken"`
}

type SecurityKeysConfig struct {
	DaysToExpireKeys            int    `json:"daysToExpireKeys"`
	MinutesToRefreshPrivateKeys int    `json:"minutesToRefreshPrivateKeys"`
	SavePublicKeyToFile         bool   `json:"savePublicKeyToFile"`
	FileECPPublicKey            string `json:"fileECPPublicKey"`
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
	EmailDefault string `json:"emailDefault"`
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
		config.MongoDB.Host = os.Getenv("MONGO_HOST")
		config.MongoDB.Port = os.Getenv("MONGO_PORT")
		if config.SMTPServer.Host != "" {
			config.SMTPServer.Host = os.Getenv("SMTP_HOST")
		}
		config.Nats.Url = os.Getenv("NATS_URL")
		config.Nats.ClientId += "_" + uuid.NewV4().String()

		fmt.Printf("ENVIRONMENT: production\n")
	}

	fmt.Printf("MONGO_HOST: %s\nMONGO_PORT: %s\n", config.MongoDB.Host, config.MongoDB.Port)
	// fmt.Printf("MONGO_USER: %s\nMONGO_PASSWORD: %s\n", config.MongoDB.User, config.MongoDB.Password)

	return config
}
