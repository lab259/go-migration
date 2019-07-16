package migration_test

import (
	"os"
	"path"
	"testing"

	"github.com/jamillosantos/macchiato"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

func TestMigration(t *testing.T) {
	RegisterFailHandler(Fail)
	description := "Migration Test Suite"

	if os.Getenv("CI") == "" {
		macchiato.RunSpecs(t, description)
	} else {
		reporterOutputDir := "./test-results/go-migration"
		os.MkdirAll(reporterOutputDir, os.ModePerm)
		junitReporter := reporters.NewJUnitReporter(path.Join(reporterOutputDir, "results.xml"))
		macchiatoReporter := macchiato.NewReporter()
		ginkgo.RunSpecsWithCustomReporters(t, description, []ginkgo.Reporter{macchiatoReporter, junitReporter})
	}
}
