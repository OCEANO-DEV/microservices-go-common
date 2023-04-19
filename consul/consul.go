package consul

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
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

	var check_port int
	address := hostname()

	// port, err := strconv.Atoi(strings.Split(config.ListenPort, ":")[1])
	port, err := getPort(address)
	if err != nil {
		return err
	}

	if len(strings.TrimSpace(config.GrpcServer.Port)) == 0 {
		check_port = port
	} else {
		port, err = strconv.Atoi(strings.Split(config.GrpcServer.Port, ":")[1])
		if err != nil {
			return err
		}

		check_port, err = strconv.Atoi(strings.Split(config.ListenPort, ":")[1])
		if err != nil {
			return err
		}
	}

	// serviceID := config.AppName
	serviceID := fmt.Sprintf("%s-%s:%v", config.AppName, address, port)

	httpCheck := fmt.Sprintf("https://%s:%v/healthy", address, check_port)
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
		log.Printf("failed consul to register service: %s:%v ", address, port)
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

func getPort(hostname string) (int, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return 0, err
	}

	filters := filters.NewArgs()
	filters.Add("hostname", hostname)

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{Filters: filters})
	if err != nil {
		return 0, err
	}

	if len(containers) == 0 {
		return 0, nil
	}

	return int(containers[0].Ports[0].PublicPort), nil
}
