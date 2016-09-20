package metrics

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Wikia/metrics-fetcher/models"
	"fmt"
	"strings"
	"io"
	"github.com/pkg/errors"
)

func OutputMetrics(filteredMetrics []models.FilteredMetrics, writer io.Writer) error {
	log.Info("outputting metrics");
	for _, metric := range filteredMetrics {
		if (len(metric.Fields) == 0) {
			return errors.Errorf("no fields in metric %s")
		}

		tagKeysAndValues := []string{}
		for tagKey, tagValue := range metric.Tags {
			tagKeysAndValues = append(tagKeysAndValues, fmt.Sprintf("%s=%s", tagKey, tagValue))
		}
		tags := strings.Join(tagKeysAndValues, ",")

		fieldKeysAndValues := []string{}
		for fieldKey, fieldValue := range metric.Fields {
			fieldKeysAndValues = append(fieldKeysAndValues, fmt.Sprintf("%s=%s", fieldKey, fieldValue))
		}
		fields := strings.Join(fieldKeysAndValues, ",")
		fmt.Fprint(writer, fmt.Sprintf("%s,%s %s\n", "resources", tags, fields))
	}
	return nil
}