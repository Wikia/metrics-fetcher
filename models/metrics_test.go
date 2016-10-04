package models_test

import (
	"fmt"

	. "github.com/Wikia/metrics-fetcher/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Metrics", func() {
	Describe("ServiceInfo", func() {
		service := ServiceInfo{
			Host: "127.0.0.1",
			Port: 1234,
		}

		It("GetAddress() should return proper URI", func() {
			Expect(service.GetAddress()).To(Equal("http://127.0.0.1:1234/metrics"))
		})
	})

	Describe("PandoraGauge", func() {
		gauge := PandoraGauge{
			Value: []byte("12.35"),
		}

		It("ToString() Should return properly formated string", func() {
			Expect(fmt.Sprint(gauge)).To(Equal("12.350000"))
		})

		It("Parse() Should return correct float value", func() {
			Expect(gauge.Parse()).To(Equal(float64(12.35)))
		})
	})

	Describe("PandoraMeter", func() {
		meter := PandoraMeter{
			Count: 123,
		}

		It("ToString() Should return properly formated string", func() {
			Expect(fmt.Sprint(meter)).To(Equal("123"))
		})
	})

	Describe("PandoraTimer", func() {
		timer := PandoraTimer{
			Count:  123,
			P50:    45.24,
			P99:    356.23,
			M1Rate: 62.11,
		}

		It("ToString() Should return properly formated string", func() {
			Expect(fmt.Sprint(timer)).To(Equal("value: 123, P50: 45.240000, P99: 356.230000, M1_Rate: 62.110000"))
		})
	})

	Describe("FilteredMetrics", func() {
		metrics := FilteredMetrics{
			Measurement: "foo",
			Tags: map[string]string{
				"tag1": "value1",
				"tag2": "value2",
			},
			Fields: map[string]interface{}{
				"integer": 12,
				"float":   3.14,
				"string":  "some string",
			},
		}

		It("ToString() Should return properly formated string", func() {
			Expect(fmt.Sprint(metrics)).To(Equal("foo,tag1=value1,tag2=value2 float=3.14,integer=12,string=\"some string\""))
		})
	})
})
