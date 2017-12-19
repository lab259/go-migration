package migration_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/jamillosantos/migration"
)

type migrationMock struct {
	id          time.Time
	description string
	up          bool
	down        bool
	manager     migration.Manager
}

func (m *migrationMock) GetID() time.Time {
	return m.id
}

func (m *migrationMock) GetDescription() string {
	return m.description
}

func (m *migrationMock) Do() error {
	return nil
}

func (m *migrationMock) Undo() error {
	return nil
}

func (m *migrationMock) GetManager() migration.Manager {
	return m.manager
}

func (m *migrationMock) SetManager(manager migration.Manager) migration.Migration {
	m.manager = manager
	return m
}

var _ = Describe("Source Code", func() {
	Describe("List", func() {
		m1 := &migrationMock{
			id:          time.Date(2001, 1, 1, 1, 1, 1, 1, time.UTC),
			description: "GetDescription 1",
			up:          true,
			down:        true,
		}

		m2 := &migrationMock{
			id:          time.Date(2001, 1, 1, 1, 1, 1, 0, time.UTC),
			description: "GetDescription 2",
			up:          true,
			down:        false,
		}

		m3 := &migrationMock{
			id:          time.Date(2000, 1, 1, 1, 1, 1, 0, time.UTC),
			description: "GetDescription 3",
			up:          true,
			down:        false,
		}

		It("It should register migrations", func() {
			d := migration.NewCodeSource()
			d.Register(m1)
			d.Register(m2)

			list, err := d.List()
			Expect(err).To(BeNil())
			Expect(list).To(HaveLen(2))
		})

		It("It should register migrations with inverted order", func() {
			d := migration.NewCodeSource()
			d.Register(m1)
			d.Register(m2)
			d.Register(m3)

			list, err := d.List()
			Expect(err).To(BeNil())
			Expect(list).To(HaveLen(3))
			Expect(list[0].GetDescription()).To(Equal("GetDescription 3"))
			Expect(list[1].GetDescription()).To(Equal("GetDescription 2"))
			Expect(list[2].GetDescription()).To(Equal("GetDescription 1"))
		})
	})
})
