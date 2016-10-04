package metrics_test

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMetrics(t *testing.T) {
	RegisterFailHandler(Fail)
	log.SetLevel(log.ErrorLevel)
	RunSpecs(t, "Metrics Suite")
}
