package main

import (
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/cloudfoundry/python-buildpack/src/python/finalize"
	_ "github.com/cloudfoundry/python-buildpack/src/python/hooks"
	"github.com/cloudfoundry/python-buildpack/src/python/pyfinder"
	"github.com/cloudfoundry/python-buildpack/src/python/requirements"

	"github.com/cloudfoundry/libbuildpack"
)

func main() {
	logfile, err := ioutil.TempFile("", "cloudfoundry.python-buildpack.finalize")
	defer logfile.Close()
	if err != nil {
		logger := libbuildpack.NewLogger(os.Stdout)
		logger.Error("Unable to create log file: %s", err.Error())
		os.Exit(8)
	}

	stdout := io.MultiWriter(os.Stdout, logfile)
	logger := libbuildpack.NewLogger(stdout)

	buildpackDir, err := libbuildpack.GetBuildpackDir()
	if err != nil {
		logger.Error("Unable to determine buildpack directory: %s", err.Error())
		os.Exit(9)
	}

	manifest, err := libbuildpack.NewManifest(buildpackDir, logger, time.Now())
	if err != nil {
		logger.Error("Unable to load buildpack manifest: %s", err.Error())
		os.Exit(10)
	}

	stager := libbuildpack.NewStager(os.Args[1:], logger, manifest)

	if err = manifest.ApplyOverride(stager.DepsDir()); err != nil {
		logger.Error("Unable to apply override.yml files: %s", err)
		os.Exit(17)
	}

	if err := stager.SetStagingEnvironment(); err != nil {
		logger.Error("Unable to setup environment variables: %s", err.Error())
		os.Exit(11)
	}

	f := finalize.Finalizer{
		Stager:         stager,
		Manifest:       manifest,
		Log:            logger,
		Logfile:        logfile,
		Command:        &libbuildpack.Command{},
		ManagePyFinder: pyfinder.ManagePyFinder{},
		Requirements:   requirements.Reqs{},
	}

	if err := finalize.Run(&f); err != nil {
		os.Exit(12)
	}

	if err := libbuildpack.RunAfterCompile(stager); err != nil {
		logger.Error("After Compile: %s", err.Error())
		os.Exit(13)
	}

	if err := stager.SetLaunchEnvironment(); err != nil {
		logger.Error("Unable to setup launch environment: %s", err.Error())
		os.Exit(14)
	}

	stager.StagingComplete()
}
