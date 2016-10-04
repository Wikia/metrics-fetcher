package registry_test

import (
	log "github.com/Sirupsen/logrus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRegistry(t *testing.T) {
	RegisterFailHandler(Fail)
	log.SetLevel(log.ErrorLevel)
	RunSpecs(t, "Registry Suite")
}
