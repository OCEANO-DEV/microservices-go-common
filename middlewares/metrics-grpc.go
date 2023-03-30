package middlewares

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/oceano-dev/microservices-go-common/metrics"
	"github.com/oceano-dev/microservices-go-common/services"
	"google.golang.org/grpc"
)

func MetricsGRPC(service services.Metrics) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		status := http.StatusOK

		resp, err = handler(ctx, req)
		if err != nil {
			log.Println(err)
		}

		appMetric := metrics.NewHttpMetrics(info.FullMethod, "POST")
		appMetric.Started()
		appMetric.Finished()
		appMetric.StatusCode = strconv.Itoa(status)
		service.SaveHttp(appMetric)

		return resp, err
	}
}
