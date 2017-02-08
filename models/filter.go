package models

import (
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
)

const (
	filterMeter     = "meters"
	filterGauge     = "gauges"
	filterHistogram = "histograms"
	filterTimer     = "timers"
)

// Filter defines metric filters to be applied
type Filter struct {
	Group       string
	Path        string
	Measurement string
}

func (f Filter) match(key string) (bool, error) {
	match, err := regexp.MatchString(f.Path, key)
	if err != nil {
		err = errors.Wrap(err, 0)
		log.WithError(err).WithFields(log.Fields{"path": f.Path, "metric": key}).Error("Error matching filter to a metric")
		return false, err
	}

	return match, nil
}

func (f Filter) parseGauge(key string, serviceInfo ServiceInfo, metric PandoraGauge) FilteredMetrics {
	log.Debugf("Found gauge metric %s : %s", key, metric)
	finalMetric := NewFilteredMetric()
	finalMetric.Tags = map[string]string{
		"service_name": serviceInfo.Name,
		"host":         serviceInfo.Host,
		"metric_name":  key,
	}
	finalMetric.Measurement = f.Measurement
	finalMetric.Fields["value"] = metric.Parse()
	finalMetric.Fields["service_id"] = serviceInfo.ID

	return finalMetric
}

func (f Filter) parseMeter(key string, serviceInfo ServiceInfo, metric PandoraMeter) FilteredMetrics {
	log.Debugf("Found meter metric %s : %s", key, metric)
	finalMetric := NewFilteredMetric()
	finalMetric.Measurement = f.Measurement
	finalMetric.Fields["value"] = metric.Count
	finalMetric.Fields["m1_rate"] = metric.M1Rate
	finalMetric.Fields["service_id"] = serviceInfo.ID
	finalMetric.Tags = map[string]string{
		"service_name": serviceInfo.Name,
		"host":         serviceInfo.Host,
		"metric_name":  key,
	}

	return finalMetric
}

func (f Filter) parseTimer(key string, serviceInfo ServiceInfo, metric PandoraTimer) FilteredMetrics {
	log.Debugf("Found timer metric %s : %s", key, metric)
	finalMetric := NewFilteredMetric()
	finalMetric.Measurement = f.Measurement
	finalMetric.Fields["value"] = metric.Count
	finalMetric.Fields["m1_rate"] = metric.M1Rate
	finalMetric.Fields["p50"] = metric.P50
	finalMetric.Fields["p99"] = metric.P99
	finalMetric.Fields["service_id"] = serviceInfo.ID
	finalMetric.Tags = map[string]string{
		"service_name": serviceInfo.Name,
		"host":         serviceInfo.Host,
		"metric_name":  key,
	}

	return finalMetric
}

func (f Filter) averageGauges(key string, serviceName string, gauges []PandoraGauge) FilteredMetrics {
	finalMetric := NewFilteredMetric()

	if len(gauges) == 0 {
		return finalMetric
	}

	finalMetric.Measurement = "metric_graphs"
	finalMetric.Tags = map[string]string{
		"service_name": serviceName,
		"metric_name":  key,
	}

	var sum, min, max float64
	for i, item := range gauges {
		value := item.Parse()

		if i == 0 {
			max = value
			min = value
			sum = value
			continue
		}

		if value > max {
			max = value
		}
		if value < min {
			min = value
		}
		sum = sum + value
	}

	finalMetric.Fields["count"] = len(gauges)
	finalMetric.Fields["min"] = min
	finalMetric.Fields["max"] = max
	finalMetric.Fields["sum"] = sum
	finalMetric.Fields["avg"] = sum / float64(len(gauges))

	return finalMetric
}

func (f Filter) averageMeters(key string, serviceName string, meters []PandoraMeter) FilteredMetrics {
	finalMetric := NewFilteredMetric()

	if len(meters) == 0 {
		return finalMetric
	}

	finalMetric.Measurement = "metric_graphs"
	finalMetric.Tags = map[string]string{
		"service_name": serviceName,
		"metric_name":  key,
	}

	var sum uint64
	var m1RateSum float64
	for _, meter := range meters {
		sum = sum + meter.Count
		m1RateSum = m1RateSum + meter.M1Rate
	}

	finalMetric.Fields["count"] = len(meters)
	finalMetric.Fields["m1_rate"] = m1RateSum
	finalMetric.Fields["value"] = sum

	return finalMetric
}

func (f Filter) averageTimers(key string, serviceName string, timers []PandoraTimer) FilteredMetrics {
	finalMetric := NewFilteredMetric()

	if len(timers) == 0 {
		return finalMetric
	}

	finalMetric.Measurement = "metric_graphs"
	finalMetric.Tags = map[string]string{
		"service_name": serviceName,
		"metric_name":  key,
	}

	var sum uint64
	var m1Min, m1Max, m1Avg, p50Min, p50Max, p50Avg, p99Min, p99Max, p99Avg float64
	for i, timer := range timers {
		sum = sum + timer.Count

		if i == 0 {
			m1Min = timer.M1Rate
			m1Max = timer.M1Rate
			m1Avg = timer.M1Rate

			p50Min = timer.P50
			p50Max = timer.P50
			p50Avg = timer.P50

			p99Min = timer.P99
			p99Max = timer.P99
			p99Avg = timer.P99

			continue
		}

		if timer.M1Rate < m1Min {
			m1Min = timer.M1Rate
		}
		if timer.M1Rate > m1Max {
			m1Max = timer.M1Rate
		}

		if timer.P50 < p50Min {
			p50Min = timer.P50
		}
		if timer.P50 > p50Max {
			p50Max = timer.P50
		}

		if timer.P99 < p99Min {
			p99Min = timer.P99
		}
		if timer.P99 > p99Max {
			p99Max = timer.P99
		}

		m1Avg = m1Avg + timer.M1Rate
		p50Avg = p50Avg + timer.P50
		p99Avg = p99Avg + timer.P99
	}

	finalMetric.Fields["count"] = len(timers)
	finalMetric.Fields["sum"] = sum
	finalMetric.Fields["avg"] = float64(sum) / float64(len(timers))
	finalMetric.Fields["m1_min"] = m1Min
	finalMetric.Fields["m1_max"] = m1Max
	finalMetric.Fields["m1_avg"] = m1Avg / float64(len(timers))
	finalMetric.Fields["p50_min"] = p50Min
	finalMetric.Fields["p50_max"] = p50Max
	finalMetric.Fields["p50_avg"] = p50Avg / float64(len(timers))
	finalMetric.Fields["p99_min"] = p99Min
	finalMetric.Fields["p99_max"] = p99Max
	finalMetric.Fields["p99_avg"] = p99Avg / float64(len(timers))

	return finalMetric
}

// ParseSingle will parse and filter single metric
func (f Filter) ParseSingle(metrics SimpleMetrics) []FilteredMetrics {
	results := []FilteredMetrics{}
	log.Debugf("Filtering for %v", f)

	switch f.Group {
	case filterGauge:
		for k, v := range metrics.Metrics.Gauges {
			if match, _ := f.match(k); !match {
				continue
			}

			results = append(results, f.parseGauge(k, metrics.Service, v))
		}
	case filterMeter:
		for k, v := range metrics.Metrics.Meters {
			if match, _ := f.match(k); !match {
				continue
			}

			results = append(results, f.parseMeter(k, metrics.Service, v))
		}
	case filterTimer:
		for k, v := range metrics.Metrics.Timers {
			if match, _ := f.match(k); !match {
				continue
			}

			results = append(results, f.parseTimer(k, metrics.Service, v))
		}
	default:
		log.Errorf("Unknown filter group: %s", f.Group)
	}

	return results
}

// ParseMany will try to parse end group metrics
func (f Filter) ParseMany(serviceName string, metrics []SimpleMetrics) []FilteredMetrics {
	results := []FilteredMetrics{}
	log.Debugf("Groupping for %v", f)

	switch f.Group {
	case filterGauge:
		gauges := map[string][]PandoraGauge{}
		for _, metric := range metrics {
			for k, v := range metric.Metrics.Gauges {
				if match, _ := f.match(k); !match {
					continue
				}

				gauges[k] = append(gauges[k], v)
			}
		}
		for k, v := range gauges {
			results = append(results, f.averageGauges(k, serviceName, v))
		}
	case filterMeter:
		meters := map[string][]PandoraMeter{}
		for _, metric := range metrics {
			for k, v := range metric.Metrics.Meters {
				if match, _ := f.match(k); !match {
					continue
				}

				meters[k] = append(meters[k], v)
			}
		}
		for k, v := range meters {
			results = append(results, f.averageMeters(k, serviceName, v))
		}
	case filterTimer:
		timers := map[string][]PandoraTimer{}
		for _, metric := range metrics {
			for k, v := range metric.Metrics.Timers {
				if match, _ := f.match(k); !match {
					continue
				}

				timers[k] = append(timers[k], v)
			}
		}
		for k, v := range timers {
			results = append(results, f.averageTimers(k, serviceName, v))
		}
	default:
		log.Errorf("Unknown filter group: %s", f.Group)
	}

	return results
}
