package metrics

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Wikia/metrics-fetcher/models"
	"github.com/go-errors/errors"
	"github.com/influxdata/influxdb/client/v2"
)

func SendMetrics(address string, database string, retention string, username string, password string, filteredMetrics []models.FilteredMetrics, extraTags map[string]string, timestamp time.Time) error {
	if len(filteredMetrics) == 0 {
		return nil
	}

	log.WithField("db_host", address).Info("Connecting to InfluxDB")
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:      address,
		Username:  username,
		Password:  password,
		UserAgent: "metrics-fetcher",
	})

	if err != nil {
		err = errors.Wrap(err, 0)
		log.WithError(err).WithField("db_host", address).Error("Error connecting to the database")
		return err
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:        database,
		RetentionPolicy: retention,
	})

	if err != nil {
		err = errors.Wrap(err, 0)
		log.WithError(err).Error("Error creating batch")
		return err
	}

	pointsAdded := 0
	for _, metrics := range filteredMetrics {
		if len(metrics.Fields) == 0 {
			log.Warn("No fields in metric - skipping")
			continue
		}
		log.WithField("measurement", metrics.Measurement).Info("Sending metrics")

		tags := metrics.Tags
		for k, v := range extraTags {
			tags[k] = v
		}

		pt, err := client.NewPoint(metrics.Measurement, tags, metrics.Fields, timestamp)
		if err != nil {
			log.WithError(err).Error("Error adding points to a batch")
			continue
		}
		bp.AddPoint(pt)
		pointsAdded++
	}

	if pointsAdded == 0 {
		log.Warn("No points added to a batch - not sending to Influx")
		return nil
	}

	err = c.Write(bp)
	if err != nil {
		err = errors.Wrap(err, 0)
		log.WithError(err).Error("Error sending metrics to InfluxDB")
	}

	return nil
}
