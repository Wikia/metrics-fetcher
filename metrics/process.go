package metrics

import "github.com/Wikia/metrics-fetcher/models"

func Combine(serviceMetrics models.GrouppedMetrics, filters []models.Filter) ([]models.FilteredMetrics, error) {
	result := []models.FilteredMetrics{}

	for serviceName, metrics := range serviceMetrics {
		for _, filter := range filters {
			for _, metric := range metrics {
				filteredMetrics := filter.ParseSingle(metric)
				result = append(result, filteredMetrics...)
			}

			combinedMetrics := filter.ParseMany(serviceName, metrics)
			result = append(result, combinedMetrics...)
		}
	}

	return result, nil
}
