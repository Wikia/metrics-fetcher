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
			tagKeysAndValues = append(tagKeysAndValues, fmt.Sprintf("%s=%s", escapeCommasEqualSignsAndSpaces(tagKey), escapeCommasEqualSignsAndSpaces(tagValue)))
		}
		tags := strings.Join(tagKeysAndValues, ",")

		fieldKeysAndValues := []string{}
		for fieldKey, fieldValue := range metric.Fields {
			valueFormat := "%s"
			switch fieldValue.(type) {
			case string:
				valueFormat = "%q" //puts into quotes and escape quotes
			case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64:
				valueFormat = "%d"
			case float32, float64:
				valueFormat = "%g"
			case bool:
				valueFormat = "%t"
			}

			fieldKeysAndValues = append(fieldKeysAndValues, fmt.Sprintf(fmt.Sprintf("%%s=%s", valueFormat), escapeCommasEqualSignsAndSpaces(fieldKey), fieldValue))
		}
		fields := strings.Join(fieldKeysAndValues, ",")
		if (len(tags) == 0) {
			fmt.Fprint(writer, fmt.Sprintf("%s %s\n", "resources", fields))
		} else {
			fmt.Fprint(writer, fmt.Sprintf("%s,%s %s\n", "resources", tags, fields))
		}
	}
	return nil
}

func escapeCommasEqualSignsAndSpaces(value string) string {
	value = strings.Replace(value, " ", "\\ ", -1)
	value = strings.Replace(value, ",", "\\,", -1)
	value = strings.Replace(value, "=", "\\=", -1)
	return value
}