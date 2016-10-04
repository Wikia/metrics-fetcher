package registry_test

import (
	"net/http"

	"github.com/Wikia/metrics-fetcher/models"
	. "github.com/Wikia/metrics-fetcher/registry"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var appsResponse = `{"apps":[{"args":null,"backoffFactor":1.15,"backoffSeconds":1,"cmd":"python3 -m http.server 8080","constraints":[],"container":{"docker":{"image":"python:3","network":"BRIDGE","portMappings":[{"containerPort":8080,"hostPort":0,"servicePort":9000,"protocol":"tcp"},{"containerPort":161,"hostPort":0,"protocol":"udp"}]},"type":"DOCKER","volumes":[]},"cpus":0.5,"dependencies":[],"deployments":[],"disk":0.0,"env":{},"executor":"","healthChecks":[{"command":null,"gracePeriodSeconds":5,"intervalSeconds":20,"maxConsecutiveFailures":3,"path":"/","portIndex":0,"protocol":"HTTP","timeoutSeconds":20}],"id":"/fake-app","instances":2,"mem":64.0,"ports":[10000,10001],"requirePorts":false,"storeUrls":[],"tasksRunning":2,"tasksStaged":0,"upgradeStrategy":{"minimumHealthCapacity":1.0},"uris":[],"user":null,"version":"2014-09-25T02:26:59.256Z"}]}`
var appResponse = `{"app":{"args":null,"backoffFactor":1.15,"backoffSeconds":1,"cmd":"python toggle.py $PORT0","constraints":[],"container":{"docker":{"image":"python:3","network":"BRIDGE","portMappings":[{"containerPort":8080,"hostPort":0,"servicePort":9000,"protocol":"tcp"}]},"type":"DOCKER","volumes":[]},"cpus":0.2,"dependencies":[],"deployments":[],"disk":0.0,"env":{},"executor":"","healthChecks":[{"command":null,"gracePeriodSeconds":5,"intervalSeconds":10,"maxConsecutiveFailures":3,"path":"/health","portIndex":0,"protocol":"HTTP","timeoutSeconds":10}],"id":"/fake-app","instances":2,"lastTaskFailure":{"appId":"/toggle","host":"10.141.141.10","message":"Abnormal executor termination","state":"TASK_FAILED","taskId":"toggle.cc427e60-5046-11e4-9e34-56847afe9799","timestamp":"2014-09-12T23:23:41.711Z","version":"2014-09-12T23:28:21.737Z"},"mem":32.0,"ports":[10000],"requirePorts":false,"storeUrls":[],"tasks":[{"appId":"/toggle","healthCheckResults":[{"alive":true,"consecutiveFailures":0,"firstSuccess":"2014-09-13T00:20:28.101Z","lastFailure":null,"lastSuccess":"2014-09-13T00:25:07.506Z","taskId":"toggle.802df2ae-3ad4-11e4-a400-56847afe9799"}],"host":"10.141.141.10","id":"toggle.802df2ae-3ad4-11e4-a400-56847afe9799","ports":[31045],"stagedAt":"2014-09-12T23:28:28.594Z","startedAt":"2014-09-13T00:24:46.959Z","version":"2014-09-12T23:28:21.737Z"},{"appId":"/toggle","healthCheckResults":[{"alive":true,"consecutiveFailures":0,"firstSuccess":"2014-09-13T00:20:28.101Z","lastFailure":null,"lastSuccess":"2014-09-13T00:25:07.508Z","taskId":"toggle.7c99814d-3ad4-11e4-a400-56847afe9799"}],"host":"10.141.141.10","id":"toggle.7c99814d-3ad4-11e4-a400-56847afe9799","ports":[31234],"stagedAt":"2014-09-12T23:28:22.587Z","startedAt":"2014-09-13T00:24:46.965Z","version":"2014-09-12T23:28:21.737Z"}],"tasksRunning":2,"tasksStaged":0,"upgradeStrategy":{"minimumHealthCapacity":1.0},"uris":["http://downloads.mesosphere.com/misc/toggle.tgz"],"user":null,"version":"2014-09-12T23:28:21.737Z"}}`
var tasksResponse = `{"tasks":[{"id":"1"},{"id":"2"}]}`

var _ = Describe("Marathon", func() {
	var marathon *MarathonRegistry
	var server *ghttp.Server

	BeforeEach(func() {
		var err error

		server = ghttp.NewServer()
		server.AllowUnhandledRequests = true
		server.UnhandledRequestStatusCode = http.StatusNotFound
		marathon, err = NewMarathonRegistry(server.URL(), 1, nil)

		Expect(err).NotTo(HaveOccurred())
	})
	AfterEach(func() {
		server.Close()
	})

	Describe("GetServices()", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v2/apps"),
					ghttp.RespondWith(http.StatusOK, appsResponse),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v2/apps/fake-app"),
					ghttp.RespondWith(http.StatusOK, appResponse),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v2/apps/fake-app/tasks"),
					ghttp.RespondWith(http.StatusOK, tasksResponse),
				),
			)
		})

		It("Should return list of valid services", func() {
			services, err := marathon.GetServices("test")
			expectedServices := []models.ServiceInfo{
				{
					Name: "/toggle",
					ID:   "toggle.802df2ae-3ad4-11e4-a400-56847afe9799",
					Host: "10.141.141.10",
					Port: 31045,
				},
				{
					Name: "/toggle",
					ID:   "toggle.7c99814d-3ad4-11e4-a400-56847afe9799",
					Host: "10.141.141.10",
					Port: 31234,
				},
			}

			Expect(err).NotTo(HaveOccurred())
			Expect(services).To(HaveLen(2))
			Expect(services).To(ConsistOf(expectedServices))
		})
	})
})
