package metrics

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Wikia/metrics-fetcher/registry"
	"github.com/parnurzeal/gorequest"
)

type SimpleMetrics struct {
	Service registry.ServiceInfo
	Metrics map[string]interface{}
}

func GatherServiceMetrics(services []registry.ServiceInfo, queueSize int, maxWorkers int) map[string][]SimpleMetrics {
	log.Infof("Starting metrics fetching: %d services", len(services))

	gatherQueue := make(chan registry.ServiceInfo, queueSize)
	gatherResults := make(chan SimpleMetrics)
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

	metrics := make(map[string][]SimpleMetrics)

	for i := 0; i < taskNum; i++ {
		metric := <-gatherResults
		metrics[metric.Service.Name] = append(metrics[metric.Service.Name], metric)
	}

	return metrics
}

func getServiceMetrics(queue <-chan registry.ServiceInfo, results chan<- SimpleMetrics) {
	for serviceInfo := range queue {
		metric := SimpleMetrics{Service: serviceInfo}

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
