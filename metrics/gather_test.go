package metrics_test

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"

	. "github.com/Wikia/metrics-fetcher/metrics"
	"github.com/Wikia/metrics-fetcher/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var sampleJson = `{"gauges":{"io.dropwizard.jetty.MutableServletContextHandler.percent-4xx-15m":{"value":0.4038381461777269},"io.dropwizard.jetty.MutableServletContextHandler.percent-4xx-1m":{"value":0.0005553371538399754},"io.dropwizard.jetty.MutableServletContextHandler.percent-4xx-5m":{"value":0.21669558216401322},"jvm.threads.waiting.count":{"value":6}},"counters":{"io.dropwizard.jetty.MutableServletContextHandler.active-dispatches":{"count":0}},"meters":{"ch.qos.logback.core.Appender.all":{"count":5834,"m15_rate":0.0019375608220075684,"m1_rate":7.689551148484769E-16,"m5_rate":0.00002883083113130127,"mean_rate":0.016637880040457177,"units":"events/second"}},"timers":{"com.wikia.exampleservice.resources.HelloWorldResource.getHelloWorld":{"count":3,"max":0.025767076000000003,"mean":0.0008590653402105778,"min":0.00045660100000000006,"p50":0.0012157510000000002,"p75":0.0012157510000000002,"p95":0.0012157510000000002,"p98":0.0012157510000000002,"p99":0.0012157510000000002,"p999":0.0012157510000000002,"stddev":0.0007564833661840576,"m15_rate":5.87239588476656E-13,"m1_rate":1.280147540491966E-147,"m5_rate":6.14411021768334E-32,"mean_rate":0.00000855566917668053,"duration_units":"seconds","rate_units":"calls/second"}}}`

func rawJsonToFloat(raw json.RawMessage) (float64, error) {
	var result float64
	err := json.Unmarshal(raw, &result)
	return result, err
}

var _ = Describe("Gather", func() {
	var server *ghttp.Server
	var serverHost, serverPort string
	var serverPortInt int64

	BeforeEach(func() {
		server = ghttp.NewServer()
		serverHost, serverPort, _ = net.SplitHostPort(server.Addr())
		serverPortInt, _ = strconv.ParseInt(serverPort, 10, 64)
	})
	AfterEach(func() {
		server.Close()
	})

	Describe("GatherServiceMetrics()", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.RespondWith(http.StatusOK, sampleJson))
		})

		Context("With mocked json response custom filters", func() {
			It("should be serialized into proper structs", func() {
				services := []models.ServiceInfo{
					{
						Name: "test-service",
						ID:   "1234",
						Host: serverHost,
						Port: serverPortInt,
					},
				}

				metrics := GatherServiceMetrics(services, 5)

				Expect(metrics).To(HaveKey("test-service"))
				Expect(metrics["test-service"]).To(HaveLen(1))

				metric := metrics["test-service"][0]

				Expect(metric.Service).To(Equal(services[0]))
				Expect(metric.Metrics.Gauges).To(HaveKey("io.dropwizard.jetty.MutableServletContextHandler.percent-4xx-15m"))
				gauge, err := rawJsonToFloat(metric.Metrics.Gauges["io.dropwizard.jetty.MutableServletContextHandler.percent-4xx-15m"].Value)
				Expect(err).NotTo(HaveOccurred())
				Expect(gauge).To(Equal(0.4038381461777269))

				Expect(metric.Metrics.Timers).To(HaveKey("com.wikia.exampleservice.resources.HelloWorldResource.getHelloWorld"))
				timer := metric.Metrics.Timers["com.wikia.exampleservice.resources.HelloWorldResource.getHelloWorld"]
				Expect(timer.Count).To(Equal(uint64(3)))
				Expect(timer.M1Rate).To(Equal(1.280147540491966E-147))
				Expect(timer.P50).To(Equal(0.0012157510000000002))
				Expect(timer.P99).To(Equal(0.0012157510000000002))
			})
		})
	})
})
