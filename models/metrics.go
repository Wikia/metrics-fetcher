package models

import (
	"encoding/json"
	"fmt"
)

type ServiceInfo struct {
	Name string
	ID   string
	Host string
	Port int
}

func (s ServiceInfo) GetAddress() string {
	return fmt.Sprintf("http://%s:%d/metrics", s.Host, s.Port)
}

type SimpleMetrics struct {
	Service ServiceInfo
	Metrics PandoraMetrics
}

type PandoraGauge struct {
	Value     json.RawMessage
	valueType string
}

func (pg PandoraGauge) String() string {
	switch pg.valueType {
	case "float64":
		var val float64
		json.Unmarshal(pg.Value, &val)
		return fmt.Sprintf("%f", val)
	case "int64":
		var val int64
		json.Unmarshal(pg.Value, &val)
		return fmt.Sprintf("%d", val)
	}

	return fmt.Sprintf("%v", pg.Value)
}

type PandoraMeter struct {
	Count uint64
}

type PandoraTimer struct {
	Count  uint64
	P50    float64
	P99    float64
	M1Rate float64
}

type PandoraMetrics struct {
	Gauges map[string]PandoraGauge
	Meters map[string]PandoraMeter
	// Histrograms string
	// Counters string
	Timers map[string]PandoraTimer
}

type GrouppedMetrics map[string][]SimpleMetrics
type FilteredMetrics struct {
	Tags   map[string]string
	Fields map[string]interface{}
}

func NewFilteredMetric() FilteredMetrics {
	return FilteredMetrics{Tags: map[string]string{}, Fields: map[string]interface{}{}}
}

func (f FilteredMetrics) IsEmpty() bool {
	if len(f.Fields) == 0 {
		return true
	}

	return false
}
