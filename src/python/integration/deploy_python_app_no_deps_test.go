package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cutlass7 "github.com/cloudfoundry/libbuildpack/cutlass/v7"
)

var _ = Describe("deploying a web app without dependencies", func() {
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

	const skippingPipInstall = "Skipping 'pip install' since requirements.txt does not exist"

	Context("no requirements.txt or setup.py", func() {
		BeforeEach(func() {
			app = cutlass7.New(Fixtures("no_deps"))
			app.Buildpacks = []string{"python_buildpack"}
			app.SetEnv("BP_DEBUG", "1")
		})

		It("deploys", func() {
			PushAppAndConfirm(app)
			Expect(app.Stdout.String()).To(ContainSubstring(skippingPipInstall))
			Expect(app.GetBody("/gg")).To(ContainSubstring("Here is your output for /gg"))
		})
	})

	Context("with setup.py but not requirements.txt", func() {
		BeforeEach(func() {
			app = cutlass7.New(Fixtures("setup_py"))
			app.Buildpacks = []string{"python_buildpack"}
			app.SetEnv("BP_DEBUG", "1")
		})

		It("deploys", func() {
			PushAppAndConfirm(app)
			Expect(app.Stdout.String()).ToNot(ContainSubstring(skippingPipInstall))
			Expect(app.GetBody("/")).To(ContainSubstring("Knock Knock. Who is there?"))
		})
	})
})
