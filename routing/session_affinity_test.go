package routing

import (
	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/generator"
	"github.com/cloudfoundry-incubator/cf-test-helpers/helpers"
	"github.com/cloudfoundry/cf-acceptance-tests/helpers/assets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Session Affinity", func() {
	Context("when an app sets a JSESSIONID cookie", func() {
		var (
			// helloWorldAsset = assets.NewAssets().HelloWorld
			stickyAsset = assets.NewAssets().HelloRouting

			appName string
		)
		BeforeEach(func() {
			appName = PushApp(stickyAsset)
		})

		FContext("when an app has two instances", func() {
			BeforeEach(func() {
				ScaleAppInstances(appName, 2)
			})

			It("routes every request to the same instance", func() {
				Eventually(func() string {
					return helpers.CurlAppRoot(appName)
				}, DEFAULT_TIMEOUT).Should(ContainSubstring("hello world"))
			})
		})
		Context("when two apps share the same domain but have different context paths", func() {
		})
	})

	Context("when an app does not set a JSESSIONID cookie", func() {
		It("routes requests round robin to all instances", func() {
		})
	})

	Context("when there is one apwhen an app has a route service boundp", func() {
		var (
			appRoute                 string
			routeServiceName         string
			golangAsset              = assets.NewAssets().Golang
			loggingRouteServiceAsset = assets.NewAssets().LoggingRouteServiceZip
		)

		BeforeEach(func() {
			// push app
			appName := PushApp(golangAsset)

			routeServiceName = PushApp(loggingRouteServiceAsset)
			// push routing service

			// get app info
			appIp, appPort := GetAppInfo(appName)

			// associate routing service with app
			appRoute = generator.RandomName()
			RegisterRoute(appRoute, appIp, appPort, routeServiceName)
		})

		It("a request to the app is routed through the route service", func() {
			Eventually(func() string {
				return helpers.CurlAppRoot(appRoute)
			}, DEFAULT_TIMEOUT).Should(ContainSubstring("go, world"))

			Eventually(func() *Session {
				logs := cf.Cf("logs", "--recent", routeServiceName)
				Expect(logs.Wait(DEFAULT_TIMEOUT)).To(Exit(0))
				return logs
			}, DEFAULT_TIMEOUT).Should(Say("Response Body: go, world"))
		})
	})
})
