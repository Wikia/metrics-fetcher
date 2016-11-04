package models_test

import (
	. "github.com/Wikia/metrics-fetcher/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filter", func() {
	metrics := []SimpleMetrics{
		{
			Service: ServiceInfo{
				Name: "test-service",
				ID:   "123-45-67-89",
				Host: "localhost",
				Port: 1234,
			},
			Metrics: PandoraMetrics{
				Gauges: map[string]PandoraGauge{
					"some.very.custom_Path": {
						Value: []byte("1234"),
					},
					"some.very.custom_Path2": {
						Value: []byte("7532"),
					},
					"some_prefix_metric-sdf_34t_4hh2": {
						Value: []byte("683"),
					},
					"4some.very.custom_Path": {
						Value: []byte("895"),
					},
				},
				Meters: map[string]PandoraMeter{
					"some.very.custom_Path": {
						Count:  123,
						M1Rate: 2.0,
					},
					"some.very.custom_Path2": {
						Count:  73,
						M1Rate: 3.0,
					},
					"some_prefix_metric-sdf_34t_4hh2": {
						Count:  87,
						M1Rate: 4.0,
					},
					"4some.very.custom_Path": {
						Count:  1,
						M1Rate: 5.0,
					},
				},
				Timers: map[string]PandoraTimer{
					"timer_custom_path": {
						Count:  12,
						P50:    1.12,
						P99:    2.33,
						M1Rate: 3.14,
					},
				},
			},
		},
		{
			Service: ServiceInfo{
				Name: "test-service2",
				ID:   "456-22-11-11",
				Host: "localhost2",
				Port: 1234,
			},
			Metrics: PandoraMetrics{
				Gauges: map[string]PandoraGauge{
					"some.very.custom_Path": {
						Value: []byte("542"),
					},
					"some.very.custom_Path2": {
						Value: []byte("2235"),
					},
					"some_prefix_metric-sdf_34t_4hh2": {
						Value: []byte("892"),
					},
					"4some.very.custom_Path": {
						Value: []byte("481"),
					},
				},
				Meters: map[string]PandoraMeter{
					"some.very.custom_Path": {
						Count: 8,
					},
					"some.very.custom_Path2": {
						Count: 2,
					},
					"some_prefix_metric-sdf_34t_4hh2": {
						Count: 21,
					},
					"4some.very.custom_Path": {
						Count: 51,
					},
				},
				Timers: map[string]PandoraTimer{
					"timer_custom_path": {
						Count:  8,
						P50:    77.31,
						P99:    32.478,
						M1Rate: 4.904,
					},
				},
			},
		},
	}
	Describe("ParseSingle()", func() {
		Context("With simple gauge matching filter", func() {
			filter := Filter{
				Group:       "gauges",
				Path:        "^some.very.custom_Path$",
				Measurement: "test-measurement",
			}

			expectedMeasurement := []FilteredMetrics{
				{
					Measurement: "test-measurement",
					Tags: map[string]string{
						"service_name": "test-service",
						"host":         "localhost",
						"metric_name":  "some.very.custom_Path",
					},
					Fields: map[string]interface{}{
						"value":      float64(1234),
						"service_id": "123-45-67-89",
					},
				},
			}

			It("Should return correct metric filtered out", func() {
				result := filter.ParseSingle(metrics[0])

				Expect(result).To(HaveLen(1))
				Expect(result).To(ConsistOf(expectedMeasurement))
			})
		})

		Context("With simple meter matching filter", func() {
			filter := Filter{
				Group:       "meters",
				Path:        "^some.very.custom_Path$",
				Measurement: "test-measurement",
			}

			expectedMeasurement := []FilteredMetrics{
				{
					Measurement: "test-measurement",
					Tags: map[string]string{
						"service_name": "test-service",
						"host":         "localhost",
						"metric_name":  "some.very.custom_Path",
					},
					Fields: map[string]interface{}{
						"value":      uint64(123),
						"m1_rate":    float64(2.0),
						"service_id": "123-45-67-89",
					},
				},
			}

			It("Should return correct metric filtered out", func() {
				result := filter.ParseSingle(metrics[0])

				Expect(result).To(HaveLen(1))
				Expect(result).To(ConsistOf(expectedMeasurement))
			})
		})

		Context("With simple timer matching filter", func() {
			filter := Filter{
				Group:       "timers",
				Path:        "^timer_custom_path$",
				Measurement: "test-measurement",
			}

			expectedMeasurement := []FilteredMetrics{
				{
					Measurement: "test-measurement",
					Tags: map[string]string{
						"service_name": "test-service",
						"host":         "localhost",
						"metric_name":  "timer_custom_path",
					},
					Fields: map[string]interface{}{
						"value":      uint64(12),
						"p50":        float64(1.12),
						"p99":        float64(2.33),
						"m1_rate":    float64(3.14),
						"service_id": "123-45-67-89",
					},
				},
			}

			It("Should return correct metric filtered out", func() {
				result := filter.ParseSingle(metrics[0])

				Expect(result).To(HaveLen(1))
				Expect(result).To(ConsistOf(expectedMeasurement))
			})
		})

		Context("With pattern gauge matching filter (prefix)", func() {
			filter := Filter{
				Group:       "gauges",
				Path:        "^some_prefix_metric",
				Measurement: "test-measurement",
			}

			expectedMeasurement := []FilteredMetrics{
				{
					Measurement: "test-measurement",
					Tags: map[string]string{
						"service_name": "test-service",
						"host":         "localhost",
						"metric_name":  "some_prefix_metric-sdf_34t_4hh2",
					},
					Fields: map[string]interface{}{
						"value":      float64(683),
						"service_id": "123-45-67-89",
					},
				},
			}

			It("Should return correct metric filtered out", func() {
				result := filter.ParseSingle(metrics[0])

				Expect(result).To(HaveLen(1))
				Expect(result).To(ConsistOf(expectedMeasurement))
			})
		})

		Context("With pattern gauge matching filter (mid substring)", func() {
			filter := Filter{
				Group:       "gauges",
				Path:        "very.custom",
				Measurement: "test-measurement",
			}

			expectedMeasurement := []FilteredMetrics{
				{
					Measurement: "test-measurement",
					Tags: map[string]string{
						"service_name": "test-service",
						"host":         "localhost",
						"metric_name":  "some.very.custom_Path",
					},
					Fields: map[string]interface{}{
						"value":      float64(1234),
						"service_id": "123-45-67-89",
					},
				},
				{
					Measurement: "test-measurement",
					Tags: map[string]string{
						"service_name": "test-service",
						"host":         "localhost",
						"metric_name":  "some.very.custom_Path2",
					},
					Fields: map[string]interface{}{
						"value":      float64(7532),
						"service_id": "123-45-67-89",
					},
				},
				{
					Measurement: "test-measurement",
					Tags: map[string]string{
						"service_name": "test-service",
						"host":         "localhost",
						"metric_name":  "4some.very.custom_Path",
					},
					Fields: map[string]interface{}{
						"value":      float64(895),
						"service_id": "123-45-67-89",
					},
				},
			}

			It("Should return correct metric filtered out", func() {
				result := filter.ParseSingle(metrics[0])

				Expect(result).To(HaveLen(3))
				Expect(result).To(ConsistOf(expectedMeasurement))
			})
		})
	})

	Describe("ParseMany()", func() {
		Context("With simple gauge matching filter", func() {
			filter := Filter{
				Group:       "gauges",
				Path:        "^some.very.custom_Path$",
				Measurement: "test-measurement",
			}

			expectedMeasurement := []FilteredMetrics{
				{
					Measurement: "metric_graphs",
					Tags: map[string]string{
						"service_name": "test-service",
						"metric_name":  "some.very.custom_Path",
					},
					Fields: map[string]interface{}{
						"avg":   float64(888),
						"min":   float64(542),
						"max":   float64(1234),
						"sum":   float64(1776),
						"count": 2,
					},
				},
			}

			It("Should return correct grouped metric filtered out", func() {
				result := filter.ParseMany("test-service", metrics)

				Expect(result).To(HaveLen(1))
				Expect(result).To(ConsistOf(expectedMeasurement))
			})
		})

		Context("With simple meter matching filter", func() {
			filter := Filter{
				Group:       "meters",
				Path:        "^some.very.custom_Path$",
				Measurement: "test-measurement",
			}

			expectedMeasurement := []FilteredMetrics{
				{
					Measurement: "metric_graphs",
					Tags: map[string]string{
						"service_name": "test-service",
						"metric_name":  "some.very.custom_Path",
					},
					Fields: map[string]interface{}{
						"value":   uint64(131),
						"m1_rate": float64(2.0),
						"count":   2,
					},
				},
			}

			It("Should return correct grouped metric filtered out", func() {
				result := filter.ParseMany("test-service", metrics)

				Expect(result).To(HaveLen(1))
				Expect(result).To(ConsistOf(expectedMeasurement))
			})
		})

		Context("With simple timer matching filter", func() {
			filter := Filter{
				Group:       "timers",
				Path:        "^timer_custom_path$",
				Measurement: "test-measurement",
			}

			expectedMeasurement := []FilteredMetrics{
				{
					Measurement: "metric_graphs",
					Tags: map[string]string{
						"service_name": "test-service",
						"metric_name":  "timer_custom_path",
					},
					Fields: map[string]interface{}{
						"p99_min": float64(2.33),
						"p99_max": float64(32.478),
						"p99_avg": float64(17.404),
						"sum":     uint64(20),
						"avg":     float64(10),
						"p50_min": float64(1.12),
						"m1_avg":  float64(4.022),
						"p50_max": float64(77.31),
						"p50_avg": float64(39.215),
						"m1_min":  float64(3.14),
						"m1_max":  float64(4.904),
						"count":   2,
					},
				},
			}

			It("Should return correct grouped metric filtered out", func() {
				result := filter.ParseMany("test-service", metrics)

				Expect(result).To(HaveLen(1))
				Expect(result).To(ConsistOf(expectedMeasurement))
			})
		})
	})
})
