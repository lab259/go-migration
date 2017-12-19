package migration_test

import (
	"testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/jamillosantos/macchiato"
)

func TestMigration(t *testing.T) {
	RegisterFailHandler(Fail)
	macchiato.RunSpecs(t, "Migration Test Suite")
}
