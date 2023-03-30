package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/oceano-dev/microservices-go-common/metrics"
	"github.com/oceano-dev/microservices-go-common/services"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func MetricsGRPC(service services.Metrics) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		fmt.Println(req)

		resp, err = handler(ctx, req)
		var status = http.StatusOK
		if err != nil {
			fmt.Println(err)
		}
		appMetric := metrics.NewHttpMetrics("req.URL.Path", "req.Method")
		appMetric.Started()
		appMetric.Finished()
		appMetric.StatusCode = strconv.Itoa(status)
		service.SaveHttp(appMetric)

		return resp, err
	}
}

func MetricsGRPCHandler() http.HandlerFunc {
	handler := promhttp.Handler()

	return func(response http.ResponseWriter, request *http.Request) {
		handler.ServeHTTP(response, request)
	}
}
