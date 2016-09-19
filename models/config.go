package models

import (
	"regexp"

	log "github.com/Sirupsen/logrus"
)

const (
	filterMeter     = "meters"
	filterGauge     = "gauges"
	filterHistogram = "histograms"
	filterTimer     = "timers"
)

type Filter struct {
	Group string
	Path  string
	Type  string
}

func (f Filter) Parse(metrics PandoraMetrics) FilteredMetrics {
	finalMetric := NewFilteredMetric()
	log.Debugf("Filtering for %v", f)
	switch f.Group {
	case filterGauge:
		for k, v := range metrics.Gauges {
			match, err := regexp.MatchString(f.Path, k)
			if err != nil {
				log.WithError(err).Error("Error matching metric")
				continue
			}

			if !match {
				continue
			}

			v.valueType = f.Type
			log.Debugf("Found metric %s : %s", k, v)
			finalMetric.Fields[k] = v
		}
	}

	return finalMetric
}
