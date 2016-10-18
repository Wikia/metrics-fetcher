package metrics_test

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/Wikia/metrics-fetcher/metrics"
	"github.com/Wikia/metrics-fetcher/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Push", func() {
	var server *ghttp.Server

	testUsername := "foo"
	testPassword := "secret"

	BeforeEach(func() {
		server = ghttp.NewServer()
	})

	AfterEach(func() {
		//shut down the server between tests
		server.Close()
	})

	Describe("SendMetrics()", func() {
		timestamp := time.Now().UTC()

		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyBasicAuth(testUsername, testPassword),
					ghttp.VerifyRequest("POST", "/write", "consistency=&db=services&precision=ns&rp=default"),
					ghttp.VerifyBody([]byte(fmt.Sprintf("test-measurement,host=localhost,metric_name=test_metric,service_name=test-service service_id=\"1234-5678-90\",value=123.45566 %d\nmetric_graphs,metric_name=test_metric,service_name=test-service max=2568.4762,med=733.68,min=12.345 %d\n", timestamp.UnixNano(), timestamp.UnixNano()))),
					ghttp.RespondWith(http.StatusOK, "OK"),
				),
			)
		})
		metrics := []models.FilteredMetrics{
			{
				Measurement: "test-measurement",
				Tags: map[string]string{
					"service_name": "test-service",
					"host":         "localhost",
					"metric_name":  "test_metric",
				},
				Fields: map[string]interface{}{
					"value":      123.45566,
					"service_id": "1234-5678-90",
				},
			},
			{
				Measurement: "metric_graphs",
				Tags: map[string]string{
					"service_name": "test-service",
					"metric_name":  "test_metric",
				},
				Fields: map[string]interface{}{
					"min": 12.345,
					"max": 2568.4762,
					"med": 733.68,
				},
			},
		}

		It("Should receive all metrics", func() {
			err := SendMetrics(server.URL(), "services", "default", testUsername, testPassword, metrics, timestamp)

			Expect(err).NotTo(HaveOccurred())
			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})
	})
})
