package metrics

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Wikia/metrics-fetcher/models"
	"github.com/go-errors/errors"
	"github.com/influxdata/influxdb/client/v2"
)

func SendMetrics(address string, username string, password string, grouppedMetrics models.GrouppedMetrics) error {
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

	for serviceName, metrics := range grouppedMetrics {
		bp, err := client.NewBatchPoints(client.BatchPointsConfig{
			Database:        "services",
			RetentionPolicy: "default",
		})

		if err != nil {
			err = errors.Wrap(err, 0)
			log.WithError(err).Error("Error creating batch")
			continue
		}

		log.WithField("service_name", serviceName).Info("Sending metrics")
		tags := map[string]string{"service_name": serviceName}
		pointsNum := 0
		for _, metric := range metrics {
			fields := map[string]interface{}{}
			tags["hostname"] = metric.Service.Host
			pt, err := client.NewPoint("service_stats", tags, fields, time.Now())
			if err != nil {
				log.WithError(err).Error("Error adding points to a batch")
				continue
			}
			pointsNum++
			bp.AddPoint(pt)
		}

		if pointsNum == 0 {
			log.WithField("service_name", serviceName).Info("No points to send - skipping")
			continue
		}

		err = c.Write(bp)
		if err != nil {
			err = errors.Wrap(err, 0)
			log.WithError(err).Error("Error sending metrics to InfluxDB")
		}
	}

	return nil
}
