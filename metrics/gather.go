package metrics

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Wikia/metrics-fetcher/models"
	"github.com/parnurzeal/gorequest"
)

func GatherServiceMetrics(services []models.ServiceInfo, queueSize int, maxWorkers int) map[string][]models.SimpleMetrics {
	log.Infof("Starting metrics fetching: %d services", len(services))

	gatherQueue := make(chan models.ServiceInfo, queueSize)
	gatherResults := make(chan models.SimpleMetrics)
	defer close(gatherResults)

	taskNum := 0

	log.Debugf("Starting workers for %d jobs", len(services))

	for i := 0; i < maxWorkers; i++ {
		go getServiceMetrics(gatherQueue, gatherResults)
	}

	for i, service := range services {
		log.Debugf("Queing service '%s' (%d)", service.ID, i+1)
		gatherQueue <- service
		taskNum++
	}
	close(gatherQueue)

	metrics := make(map[string][]models.SimpleMetrics)

	for i := 0; i < taskNum; i++ {
		metric := <-gatherResults
		metrics[metric.Service.Name] = append(metrics[metric.Service.Name], metric)
	}

	return metrics
}

func getServiceMetrics(queue <-chan models.ServiceInfo, results chan<- models.SimpleMetrics) {
	for serviceInfo := range queue {
		metric := models.SimpleMetrics{Service: serviceInfo}

		log.WithFields(log.Fields{"task_id": metric.Service.ID, "uri": metric.Service.GetAddress()}).Info("Fetching metrics for service")

		resp, _, err := gorequest.New().Get(metric.Service.GetAddress()).EndStruct(&metric.Metrics)

		if len(err) != 0 {
			log.WithFields(log.Fields{"task_id": metric.Service.ID, "uri": metric.Service.GetAddress(), "errors": err}).Error("Error fetching metrics")
			results <- metric
			continue
		}

		if resp.StatusCode != 200 {
			log.WithFields(log.Fields{"task_id": metric.Service.ID, "uri": metric.Service.GetAddress()}).Error("Response status != 200")
			results <- metric
			continue
		}

		results <- metric
	}
}
