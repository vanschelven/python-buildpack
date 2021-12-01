package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cutlass7 "github.com/cloudfoundry/libbuildpack/cutlass/v7"
)

var _ = Describe("deploying a flask web app", func() {

	var app *cutlass7.App

	BeforeEach(func() {
		if isSerialTest {
			Skip("Skipping parallel tests")
		}
	})

	AfterEach(func() {
		if app != nil {
			app.Destroy()
		}
		app = nil
	})

	Context("start command is specified in manifest.yml", func() {
		BeforeEach(func() {
			app = cutlass7.New(Fixtures("flask_no_procfile"))
		})

		It("deploys", func() {
			PushAppAndConfirm(app)
			Expect(app.GetBody("/")).To(ContainSubstring("I was started without a Procfile"))
		})
	})
})
