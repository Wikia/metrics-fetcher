package registry

import (
	"net/url"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/Wikia/metrics-fetcher/models"
	marathon "github.com/gambol99/go-marathon"
	"github.com/go-errors/errors"
)

type MarathonRegistry struct {
	client    marathon.Marathon
	MaxWorker int
	MaxQueue  int
}

func NewMarathonRegistry(host string, queueSize int, numWorkers int) (*MarathonRegistry, error) {
	config := marathon.NewDefaultConfig()
	config.URL = host

	log.Debug("Configuring Marathon Client with host: ", host)

	marathonClient, err := marathon.NewClient(config)

	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return &MarathonRegistry{
		client:    marathonClient,
		MaxWorker: numWorkers,
		MaxQueue:  queueSize,
	}, nil
}

func fetchServiceTasks(client marathon.Marathon, queue <-chan string, results chan<- models.ServiceInfo, finish chan<- bool) {
	for appID := range queue {
		log.WithField("app_id", appID).Debug("Fetching tasks")
		details, err := client.Application(appID)
		if err != nil {
			log.WithError(err).WithField("app_id", appID).Error("Error getting app details")
			finish <- false
			continue
		}

		for _, task := range details.Tasks {
			log.WithField("app_id", appID).Debug("Adding task: ", task.ID)
			if len(task.Ports) == 0 {
				log.WithField("app_id", appID).Warn("Service has no ports defined: skipping")
				finish <- false
				continue
			}

			results <- models.ServiceInfo{
				Name: task.AppID,
				ID:   task.ID,
				Host: task.Host,
				Port: task.Ports[len(task.Ports)-1],
			}
		}
		log.WithField("app_id", appID).Debug("Finished adding tasks")
		finish <- true
	}
}

func (c MarathonRegistry) GetServices(label string) ([]models.ServiceInfo, error) {
	v := url.Values{}
	v.Set("label", label)

	apps, err := c.client.Applications(v)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	log.Infof("Fetched %d apps with label '%s'", len(apps.Apps), label)

	queue := make(chan string, c.MaxQueue)
	results := make(chan models.ServiceInfo)
	defer close(results)

	taskNum := 0
	processed := make(chan bool, len(apps.Apps))

	log.Debugf("Starting workers for %d jobs", taskNum-1)

	for i := 0; i < c.MaxWorker; i++ {
		go fetchServiceTasks(c.client, queue, results, processed)
	}

	for i, app := range apps.Apps {
		log.Debugf("Found application '%s' (%d)", app.ID, i+1)
		queue <- app.ID
		taskNum++
	}
	close(queue)

	var serviceInfos []models.ServiceInfo
FETCHLOOP:
	for {
		select {
		case serviceInfo := <-results:
			serviceInfos = append(serviceInfos, serviceInfo)
		case <-processed:
			log.Debug("Worker finished, left: ", taskNum)
			taskNum--
			if taskNum <= 0 {
				break FETCHLOOP
			}
		case <-time.After(time.Second * 3):
			break FETCHLOOP
		}
	}

	return serviceInfos, nil
}
