package models

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// ServiceInfo holds basic information about a service
type ServiceInfo struct {
	Name string
	ID   string
	Host string
	Port int64
}

// GetAddress returns the service address from which metrics are fetched
func (s ServiceInfo) GetAddress() string {
	return fmt.Sprintf("http://%s:%d/metrics", s.Host, s.Port)
}

// SimpleMetrics represents very simple metric for Pandora service
type SimpleMetrics struct {
	Service ServiceInfo
	Metrics PandoraMetrics
}

// PandoraGauge is the definition of gauge metric
type PandoraGauge struct {
	Value json.RawMessage
}

func (pg PandoraGauge) String() string {
	return fmt.Sprintf("%f", pg.Parse())
}

// Parse will normalize the gauge value to float64
func (pg PandoraGauge) Parse() float64 {
	var val float64
	json.Unmarshal(pg.Value, &val)

	return val
}

// PandoraMeter is the definition of meter metrics
type PandoraMeter struct {
	Count  uint64
	M1Rate float64 `json:"m1_rate"`
}

func (pm PandoraMeter) String() string {
	return fmt.Sprintf("%v", pm.Count)
}

// PandoraTimer is the definition of timer metric
type PandoraTimer struct {
	Count  uint64
	P50    float64
	P99    float64
	M1Rate float64 `json:"m1_rate"`
}

func (pt PandoraTimer) String() string {
	return fmt.Sprintf("value: %v, P50: %f, P99: %f, M1_Rate: %f", pt.Count, pt.P50, pt.P99, pt.M1Rate)
}

// PandoraMetrics defines all the metrics returned by the Pandora service
type PandoraMetrics struct {
	Gauges map[string]PandoraGauge
	Meters map[string]PandoraMeter
	// Histrograms string
	// Counters string
	Timers map[string]PandoraTimer
}

// GroupedMetrics is map of service name to an array of metrics
type GroupedMetrics map[string][]SimpleMetrics

// FilteredMetrics is the resulting metrics that is being sent to Influx database
type FilteredMetrics struct {
	Measurement string
	Tags        map[string]string
	Fields      map[string]interface{}
}

func (fm FilteredMetrics) String() string {
	tags := make([]string, len(fm.Tags))
	i := 0
	var keys []string
	for k := range fm.Tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		tags[i] = fmt.Sprintf("%s=%v", k, fm.Tags[k])
		i++
	}

	fields := make([]string, len(fm.Fields))
	i = 0
	keys = []string{}
	for k := range fm.Fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := fm.Fields[k]
		switch v.(type) {
		case string:
			fields[i] = fmt.Sprintf("%s=%s", k, strconv.Quote(v.(string)))
		default:
			fields[i] = fmt.Sprintf("%s=%v", k, v)

		}
		i++
	}
	return fmt.Sprintf("%s,%s %s", fm.Measurement, strings.Join(tags, ","), strings.Join(fields, ","))
}

// NewFilteredMetric creates new instance of FilteredMetrics
func NewFilteredMetric() FilteredMetrics {
	return FilteredMetrics{Tags: map[string]string{}, Fields: map[string]interface{}{}}
}
