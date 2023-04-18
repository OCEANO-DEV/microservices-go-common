package consul

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/oceano-dev/microservices-go-common/config"

	consul "github.com/hashicorp/consul/api"
)

func NewConsulClient(
	config *config.Config,
) (*consul.Client, error) {

	consulConfig := consul.DefaultConfig()
	consulConfig.Address = config.Consul.Host

	consulClient, err := consul.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}

	err = register(config, consulClient)
	if err != nil {
		return nil, err
	}

	return consulClient, nil
}

func register(config *config.Config, client *consul.Client) error {
	port, err := strconv.Atoi(strings.Split(config.ListenPort, ":")[1])
	if err != nil {
		return err
	}

	address := hostname()
	serviceID := config.AppName
	// serviceID := fmt.Sprintf("%s-%s:%v", config.AppName, address, port)

	httpCheck := fmt.Sprintf("https://%s:%v/healthy", address, port)
	fmt.Println(httpCheck)

	registration := &consul.AgentServiceRegistration{
		ID:      serviceID,
		Name:    config.AppName,
		Port:    port,
		Address: address,
		Check: &consul.AgentServiceCheck{
			CheckID:                        fmt.Sprintf("%s_app_status", config.AppName),
			Name:                           fmt.Sprintf("%s application status", config.AppName),
			HTTP:                           httpCheck,
			TLSSkipVerify:                  true,
			Interval:                       "10s",
			Timeout:                        "30s",
			DeregisterCriticalServiceAfter: "30m",
		},
	}

	err = client.Agent().ServiceRegister(registration)

	if err != nil {
		log.Printf("Failed consul to register service: %s:%v ", address, port)
		return err
	}

	log.Printf("successfully consul register service: %s:%v", address, port)

	return nil
}

// func deregister(config *config.Config, client *consul.Client) error {
// 	return client.Agent().ServiceDeregister(config.AppName)
// }

func hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}

	return hostname
}
