package registry

import (
	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/Wikia/metrics-fetcher/models"
	marathon "github.com/gambol99/go-marathon"
	"github.com/go-errors/errors"
	"gopkg.in/go-playground/pool.v3"
)

// MarathonRegistry is the structure used to fetch services from Marathon
type MarathonRegistry struct {
	client    marathon.Marathon
	MaxWorker uint
}

// NewMarathonRegistry creates new MarathonRegistry instance and instantiates API client
func NewMarathonRegistry(host string, numWorkers uint, client *http.Client) (*MarathonRegistry, error) {
	config := marathon.NewDefaultConfig()
	config.URL = host
	if client != nil {
		config.HTTPClient = client
	}

	log.Debug("Configuring Marathon Client with host: ", host)

	marathonClient, err := marathon.NewClient(config)

	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return &MarathonRegistry{
		client:    marathonClient,
		MaxWorker: numWorkers,
	}, nil
}

func fetchServiceTasks(client marathon.Marathon, appID string) pool.WorkFunc {
	return func(wu pool.WorkUnit) (interface{}, error) {
		log.WithField("app_id", appID).Debug("Fetching tasks")

		details, err := client.Application(appID)
		if err != nil {
			err = errors.Wrap(err, 0)
			log.WithError(err).WithField("app_id", appID).Error("Error getting app details")
			return nil, err
		}

		if wu.IsCancelled() {
			return nil, nil
		}

		result := []models.ServiceInfo{}
		for _, task := range details.Tasks {
			log.WithField("app_id", appID).Debug("Adding task: ", task.ID)
			if len(task.Ports) == 0 {
				log.WithField("app_id", appID).Warn("Service has no ports defined: skipping")
				return nil, errors.Errorf("No prort defined for service: %s", appID)
			}

			result = append(result, models.ServiceInfo{
				Name: task.AppID,
				ID:   task.ID,
				Host: task.Host,
				Port: int64(task.Ports[len(task.Ports)-1]),
			})
		}
		log.WithField("app_id", appID).Debug("Finished adding tasks")
		return result, nil
	}
}

// GetServices returns list of services with a given label
func (c MarathonRegistry) GetServices(label string) ([]models.ServiceInfo, error) {
	v := url.Values{}
	v.Set("label", label)

	apps, err := c.client.Applications(v)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	log.Infof("Fetched %d apps with label '%s'", len(apps.Apps), label)

	p := pool.NewLimited(c.MaxWorker)
	defer p.Close()

	log.Debugf("Starting workers for %d jobs", len(apps.Apps))
	batch := p.Batch()
	go func() {
		for i, app := range apps.Apps {
			log.Debugf("Found application '%s' (%d)", app.ID, i+1)
			batch.Queue(fetchServiceTasks(c.client, app.ID))
		}
		batch.QueueComplete()
	}()
	log.Debug("All tasks scheduled!")

	var serviceInfos []models.ServiceInfo
	for infos := range batch.Results() {
		if err := infos.Error(); err != nil {
			log.WithError(err).Error("Error fetching results")
			continue
		}
		log.Debug("Successfully retrieved an result")
		serviceInfos = append(serviceInfos, infos.Value().([]models.ServiceInfo)...)
	}

	return serviceInfos, nil
}
