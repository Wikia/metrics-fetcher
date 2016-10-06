// Copyright © 2016 Wikia Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"runtime"
	"time"

	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/Wikia/metrics-fetcher/metrics"
	"github.com/Wikia/metrics-fetcher/models"
	"github.com/Wikia/metrics-fetcher/registry"
	"github.com/go-errors/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	marathonHost    string
	marathonLabel   string
	influxAddress   string
	influxDB        string
	influxRetention string
	numWorkers      uint
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Gathers metrics from the services",
	Long: `First it fetches list of services from the Consul registry with a specific
tag to process. Then it calls the very last port defined on the service (assuming this is the admin port)
to fetch metrics. Then it aggregates those metrics by a service id and sends them back to Influx
For now it supports only Influx line protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		serviceRegistry, err := registry.NewMarathonRegistry(marathonHost, numWorkers, nil)
		if err != nil {
			log.Error(err)
			return
		}
		log.WithField("marathon_lable", marathonLabel).Info("Getting services for measurement")
		services, err := serviceRegistry.GetServices(marathonLabel)
		if err != nil {
			log.WithError(err).Error("Erorr getting list of services")
			return
		}

		// gathering metrics
		log.Infof("Fetching metrics from services: %d", len(services))
		grouppedMetrics := metrics.GatherServiceMetrics(services, numWorkers)

		filters := []models.Filter{}
		err = viper.UnmarshalKey("filters", &filters)

		if err != nil {
			err = errors.Wrap(err, 0)
			log.WithError(err).Error("Error loading filters from configuration")
			return
		}
		combinedMetrics, _ := metrics.Combine(grouppedMetrics, filters)
		metrics.OutputMetrics(combinedMetrics, os.Stdout)

		if len(influxAddress) != 0 {
			log.WithField("server", influxAddress).Info("Sending metrics to database")
			err = metrics.SendMetrics(influxAddress, influxDB, influxRetention, "", "", combinedMetrics, time.Now())
			if err != nil {
				log.WithError(err).Error("Error sending metrics")
				return
			}
		}
	},
}

func init() {
	fetchCmd.Flags().StringVar(&marathonHost, "marathon", "http://localhost:8080", "address of a marathon API to connect to")
	fetchCmd.Flags().StringVar(&marathonLabel, "label", "gather-metrics", "label to search services in marathon with")
	fetchCmd.Flags().StringVar(&influxAddress, "influx", "", "address of an InfluxDB server where metrics should be pushed")
	fetchCmd.Flags().StringVar(&influxDB, "database", "services", "name of the InfluxDB database")
	fetchCmd.Flags().StringVar(&influxRetention, "retention", "default", "which retention policy should we use for pushing metrics")
	fetchCmd.Flags().UintVar(&numWorkers, "workers", uint(runtime.NumCPU()*5), "how many fetcher workers to spawn")
	RootCmd.AddCommand(fetchCmd)
}
