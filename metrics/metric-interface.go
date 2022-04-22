package metrics

type IMetric interface {
	SaveClient(client *ClientMetrics) error
	SaveHttp(http *HttpMetrics)
}
