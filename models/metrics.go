package models

import (
	"encoding/json"
	"fmt"
	"strings"
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

func (pm PandoraMeter) String() string {
	return fmt.Sprintf("%v", pm.Count)
}

type PandoraTimer struct {
	Count  uint64
	P50    float64
	P99    float64
	M1Rate float64
}

func (pt PandoraTimer) String() string {
	return fmt.Sprintf("value: %v, P50: %v, P99: %v, M1_Rate: %v", pt.Count, pt.P50, pt.P99, pt.M1Rate)
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
	Measurement string
	Tags        map[string]string
	Fields      map[string]interface{}
}

func (fm FilteredMetrics) String() string {
	tags := make([]string, len(fm.Tags))
	i := 0
	for k, v := range fm.Tags {
		tags[i] = fmt.Sprintf("%s=%s", k, v)
		i++
	}

	fields := make([]string, len(fm.Fields))
	i = 0
	for k, v := range fm.Fields {
		fields[i] = fmt.Sprintf("%s=%s", k, v)
		i++
	}
	return fmt.Sprintf("%s,%s %s", fm.Measurement, strings.Join(tags, ","), strings.Join(fields, ","))
}

func NewFilteredMetric() FilteredMetrics {
	return FilteredMetrics{Tags: map[string]string{}, Fields: map[string]interface{}{}}
}
