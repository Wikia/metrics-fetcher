package metrics_test

import (
	. "github.com/Wikia/metrics-fetcher/metrics"
	"github.com/Wikia/metrics-fetcher/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Process", func() {
	Describe("CombineMetrics()", func() {
		filters := []models.Filter{
			{
				Path:        "test_metric",
				Group:       "gauges",
				Measurement: "test_measurement",
			},
		}

		Context("With one metric", func() {
			metrics := models.GrouppedMetrics{
				"test-service": []models.SimpleMetrics{
					{
						Service: models.ServiceInfo{
							Name: "test-service",
							ID:   "1234-5678-90",
							Host: "localhost",
							Port: 1234,
						},
						Metrics: models.PandoraMetrics{
							Gauges: map[string]models.PandoraGauge{
								"test_metric": {
									Value: []byte("112233.445566"),
								},
							},
						},
					},
				},
			}

			It("Should return one metric point", func() {
				measurements, err := Combine(metrics, filters)
				Expect(err).NotTo(HaveOccurred())
				Expect(measurements).To(HaveLen(2))

				expectedMetrics := []models.FilteredMetrics{
					{
						Measurement: filters[0].Measurement,
						Tags: map[string]string{
							"service_name": "test-service",
							"host":         "localhost",
							"metric_name":  "test_metric",
						},
						Fields: map[string]interface{}{
							"value":      112233.445566,
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
							"min":   112233.445566,
							"max":   112233.445566,
							"avg":   112233.445566,
							"count": 1,
						},
					},
				}
				Expect(measurements).To(ConsistOf(expectedMetrics))
			})
		})

		Context("With many metrics", func() {
			metrics := models.GrouppedMetrics{
				"test-service": []models.SimpleMetrics{
					{
						Service: models.ServiceInfo{
							Name: "test-service",
							ID:   "1234-5678-90",
							Host: "localhost1",
							Port: 1234,
						},
						Metrics: models.PandoraMetrics{
							Gauges: map[string]models.PandoraGauge{
								"test_metric": {
									Value: []byte("100.00"),
								},
							},
						},
					},
					{
						Service: models.ServiceInfo{
							Name: "test-service",
							ID:   "998877-665544-321",
							Host: "localhost2",
							Port: 1234,
						},
						Metrics: models.PandoraMetrics{
							Gauges: map[string]models.PandoraGauge{
								"test_metric": {
									Value: []byte("20.0"),
								},
							},
						},
					},
				},
			}

			It("Should return properly aggregated metrics", func() {
				measurements, err := Combine(metrics, filters)
				Expect(err).NotTo(HaveOccurred())
				Expect(measurements).To(HaveLen(3))

				expectedMetrics := []models.FilteredMetrics{
					{
						Measurement: filters[0].Measurement,
						Tags: map[string]string{
							"service_name": "test-service",
							"host":         "localhost1",
							"metric_name":  "test_metric",
						},
						Fields: map[string]interface{}{
							"value":      100.00,
							"service_id": "1234-5678-90",
						},
					},
					{
						Measurement: filters[0].Measurement,
						Tags: map[string]string{
							"service_name": "test-service",
							"host":         "localhost2",
							"metric_name":  "test_metric",
						},
						Fields: map[string]interface{}{
							"value":      20.00,
							"service_id": "998877-665544-321",
						},
					},
					{
						Measurement: "metric_graphs",
						Tags: map[string]string{
							"service_name": "test-service",
							"metric_name":  "test_metric",
						},
						Fields: map[string]interface{}{
							"min":   20.00,
							"max":   100.00,
							"avg":   60.00,
							"count": 2,
						},
					},
				}

				Expect(measurements).To(ConsistOf(expectedMetrics))
			})
		})
	})
})
