package metrics

import "github.com/Wikia/metrics-fetcher/models"

// Combine metrics and filter them according to current configuration
func Combine(serviceMetrics models.GroupedMetrics, filters []models.Filter) ([]models.FilteredMetrics, error) {
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
