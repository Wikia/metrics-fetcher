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

type Filter struct {
	Group       string
	Path        string
	Type        string
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
	metric.valueType = f.Type
	log.Debugf("Found gauge metric %s : %s", key, metric)
	finalMetric := NewFilteredMetric()
	finalMetric.Tags = map[string]string{
		"service_name": serviceInfo.Name,
		"host":         serviceInfo.Host,
		"metric_name":  key,
	}
	finalMetric.Measurement = f.Measurement
	finalMetric.Fields["value"] = metric
	finalMetric.Fields["service_id"] = serviceInfo.ID

	return finalMetric
}

func (f Filter) parseMeter(key string, serviceInfo ServiceInfo, metric PandoraMeter) FilteredMetrics {
	log.Debugf("Found meter metric %s : %s", key, metric)
	finalMetric := NewFilteredMetric()
	finalMetric.Measurement = f.Measurement
	finalMetric.Fields["value"] = metric.Count
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

func (f Filter) Parse(metrics SimpleMetrics) []FilteredMetrics {
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
