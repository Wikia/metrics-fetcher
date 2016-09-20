package metrics

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Wikia/metrics-fetcher/models"
	"github.com/go-errors/errors"
	"github.com/spf13/viper"
)

func CombineMetrics(serviceMetrics models.GrouppedMetrics) ([]models.FilteredMetrics, error) {
	result := []models.FilteredMetrics{}
	filters := []models.Filter{}
	err := viper.UnmarshalKey("filters", &filters)

	if err != nil {
		err = errors.Wrap(err, 0)
		log.WithError(err).Error("Error loading filters from configuration")
		return result, err
	}

	for _, metrics := range serviceMetrics {
		for _, filter := range filters {
			for _, metric := range metrics {
				filteredMetrics := filter.Parse(metric)
				result = append(result, filteredMetrics...)
			}
		}
	}

	return result, nil
}
