package migration_test

import (
	"errors"
	"github.com/lab259/go-migration"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("MigrationDefault", func() {
	It("should create a new instance", func() {
		now := time.Now().UTC()
		m := migration.NewMigration(now, "description 1")
		Expect(m.GetID()).To(Equal(now))
		Expect(m.GetDescription()).To(Equal("description 1"))
	})

	It("should create a new instance with one handler", func() {
		now := time.Now().UTC()

		do := func () error {
			return errors.New("this is do")
		}

		m := migration.NewMigration(now, "description 1", do)
		Expect(m.GetID()).To(Equal(now))
		Expect(m.GetDescription()).To(Equal("description 1"))
		err := m.Do()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("this is do"))
	})

	It("should create a new instance with two handlers", func() {
		now := time.Now().UTC()

		do := func () error {
			return errors.New("this is do")
		}
		undo := func () error {
			return errors.New("this is undo")
		}

		m := migration.NewMigration(now, "description 1", do, undo)
		Expect(m.GetID()).To(Equal(now))
		Expect(m.GetDescription()).To(Equal("description 1"))
		err := m.Do()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("this is do"))
		err = m.Undo()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("this is undo"))
	})

	It("should create a new instance with a manager", func() {
		manager := migration.NewDefaultManager(&nopTarget{}, migration.NewCodeSource())

		now := time.Now().UTC()
		m := migration.NewMigration(now, "description 1")
		m.SetManager(manager)
		Expect(m.GetManager()).To(Equal(manager))
	})
})
