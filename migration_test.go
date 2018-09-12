package migration_test

import (
	"github.com/jamillosantos/macchiato"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestMigration(t *testing.T) {
	RegisterFailHandler(Fail)
	macchiato.RunSpecs(t, "Migration Test Suite")
}
