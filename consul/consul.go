package consul

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/oceano-dev/microservices-go-common/config"

	consul "github.com/hashicorp/consul/api"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func NewConsulClient(
	config *config.Config,
) (*consul.Client, string, error) {

	consulConfig := consul.DefaultConfig()
	consulConfig.Address = config.Consul.Host

	consulClient, err := consul.NewClient(consulConfig)
	if err != nil {
		return nil, "", err
	}

	serviceID, err := register(config, consulClient)
	if err != nil {
		return nil, "", err
	}

	return consulClient, serviceID, nil
}

func register(config *config.Config, client *consul.Client) (string, error) {

	var check_port int
	address := hostname()
	// address := fmt.Sprintf("%s-srv", config.AppName)

	port, err := strconv.Atoi(strings.Split(config.ListenPort, ":")[1])
	if port == 0 || err != nil {
		return "", err
	}

	check_port = port

	if len(strings.TrimSpace(config.GrpcServer.Port)) > 0 {
		port, err = strconv.Atoi(strings.Split(config.GrpcServer.Port, ":")[1])
		if err != nil {
			return "", err
		}
	}

	serviceID := fmt.Sprintf("%s-%s:%v", config.AppName, address, port)

	httpCheck := fmt.Sprintf("https://%s:%v/healthy", address, check_port)
	fmt.Println(httpCheck)

	registration := &consul.AgentServiceRegistration{
		ID:      serviceID,
		Name:    config.AppName,
		Port:    port,
		Address: address,
		Check: &consul.AgentServiceCheck{
			CheckID:                        serviceID,
			Name:                           fmt.Sprintf("Service %s check", config.AppName),
			HTTP:                           httpCheck,
			TLSSkipVerify:                  true,
			Interval:                       "10s",
			Timeout:                        "30s",
			DeregisterCriticalServiceAfter: "1m",
		},
	}

	err = client.Agent().ServiceRegister(registration)

	if err != nil {
		log.Printf("failed consul to register service: %s:%v ", address, port)
		return "", err
	}

	log.Printf("successfully consul register service: %s:%v", address, port)

	return serviceID, nil
}

func hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}

	serviceName, err := getServiceNameKubernetes(hostname)
	if len(serviceName) > 0 && err == nil {
		return serviceName
	}

	return hostname
}

func getServiceNameKubernetes(podName string) (string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting Kubernetes configuration: %v\n", err)
		return "", err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating Kubernetes client: %v\n", err)
		return "", err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pod, err := clientset.CoreV1().Pods("default").Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting pod information: %v\n", err)
		return "", err
	}

	serviceName := pod.ObjectMeta.Labels["app"]
	service, err := clientset.CoreV1().Services(pod.ObjectMeta.Namespace).Get(ctx, serviceName, metav1.GetOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting information about the service: %v\n", err)
		return "", err
	}

	fmt.Printf("Pod %s belongs to service %s\n", podName, service.ObjectMeta.Name)
	return service.ObjectMeta.Name, nil
}
