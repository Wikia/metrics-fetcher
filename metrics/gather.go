package metrics

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Wikia/metrics-fetcher/models"
	"github.com/go-errors/errors"
	"github.com/parnurzeal/gorequest"
	pool "gopkg.in/go-playground/pool.v3"
)

func GatherServiceMetrics(services []models.ServiceInfo, maxWorkers uint) models.GrouppedMetrics {
	log.Infof("Starting metrics fetching: %d services", len(services))

	log.Debugf("Starting workers for %d jobs", len(services))
	p := pool.NewLimited(maxWorkers)
	defer p.Close()

	log.Debugf("Starting workers for %d jobs", len(services))
	batch := p.Batch()
	go func() {
		for i, serviceInfo := range services {
			log.Debugf("Queing service '%s' (%d)", serviceInfo.ID, i+1)
			batch.Queue(getServiceMetrics(serviceInfo))
		}
		batch.QueueComplete()
	}()
	log.Debug("All tasks scheduled!")

	metrics := make(models.GrouppedMetrics)

	for metric := range batch.Results() {
		if err := metric.Error(); err != nil {
			log.WithError(err).Error("Error fetching results")
			continue
		}
		simpleMetric := metric.Value().(models.SimpleMetrics)
		metrics[simpleMetric.Service.Name] = append(metrics[simpleMetric.Service.Name], simpleMetric)
	}

	return metrics
}

func getServiceMetrics(serviceInfo models.ServiceInfo) pool.WorkFunc {
	return func(wu pool.WorkUnit) (interface{}, error) {
		metric := models.SimpleMetrics{Service: serviceInfo}

		log.WithFields(log.Fields{"task_id": metric.Service.ID, "uri": metric.Service.GetAddress()}).Info("Fetching metrics for service")
		resp, _, err := gorequest.New().Get(metric.Service.GetAddress()).EndStruct(&metric.Metrics)

		if len(err) != 0 {
			log.WithFields(log.Fields{"task_id": metric.Service.ID, "uri": metric.Service.GetAddress(), "errors": err}).Error("Error fetching metrics")
			return nil, err[0]
		}

		if resp.StatusCode != 200 {
			log.WithFields(log.Fields{"task_id": metric.Service.ID, "uri": metric.Service.GetAddress()}).Error("Response status != 200")
			return nil, errors.Errorf("Got unexpected status from metrics endpoint: %d", resp.StatusCode)
		}

		return metric, nil
	}
}
