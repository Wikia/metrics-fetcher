package metrics

import (
	"fmt"
	"io"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/Wikia/metrics-fetcher/models"
	"github.com/pkg/errors"
)

// OutputMetrics to STDOUT
func OutputMetrics(filteredMetrics []models.FilteredMetrics, extraTags map[string]string, writer io.Writer) error {
	log.Info("outputting metrics")
	for _, metric := range filteredMetrics {
		if len(metric.Fields) == 0 {
			return errors.Errorf("no fields in metric")
		}

		tagKeysAndValues := make([]string, len(metric.Tags)+len(extraTags))
		i := 0
		for tagKey, tagValue := range metric.Tags {
			tagKeysAndValues[i] = fmt.Sprintf("%s=%s", escapeSpecialChars(tagKey), escapeSpecialChars(tagValue))
			i++
		}

		for tagKey, tagValue := range extraTags {
			tagKeysAndValues[i] = fmt.Sprintf("%s=%s", escapeSpecialChars(tagKey), escapeSpecialChars(tagValue))
			i++
		}
		tags := strings.Join(tagKeysAndValues, ",")

		fieldKeysAndValues := make([]string, len(metric.Fields))
		j := 0
		for fieldKey, fieldValue := range metric.Fields {
			var valueFormat string
			switch fieldValue.(type) {
			case string:
				valueFormat = "%q" //escapes quotes and puts into quotes
			case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64:
				valueFormat = "%d"
			case float32, float64:
				valueFormat = "%g"
			case bool:
				valueFormat = "%t"
			default:
				valueFormat = "%v"
			}

			fieldKeysAndValues[j] = fmt.Sprintf(fmt.Sprintf("%%s=%s", valueFormat), escapeSpecialChars(fieldKey), fieldValue)
			j++
		}
		fields := strings.Join(fieldKeysAndValues, ",")
		if len(tags) == 0 {
			fmt.Fprint(writer, fmt.Sprintf("%s %s\n", escapeMeasurementName(metric.Measurement), fields))
		} else {
			fmt.Fprint(writer, fmt.Sprintf("%s,%s %s\n", escapeMeasurementName(metric.Measurement), tags, fields))
		}
	}
	return nil
}

func escapeSpecialChars(value string) string {
	value = strings.Replace(value, " ", "\\ ", -1)
	value = strings.Replace(value, ",", "\\,", -1)
	value = strings.Replace(value, "=", "\\=", -1)
	return value
}

func escapeMeasurementName(value string) string {
	value = strings.Replace(value, " ", "\\ ", -1)
	value = strings.Replace(value, ",", "\\,", -1)
	return value
}
