package migration

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type migrationMock struct {
	id          time.Time
	description string
	up          bool
	down        bool
	manager     Manager
}

func (m *migrationMock) GetID() time.Time {
	return m.id
}

func (m *migrationMock) GetDescription() string {
	return m.description
}

func (m *migrationMock) Up() error {
	return nil
}

func (m *migrationMock) Down() error {
	return nil
}

func (m *migrationMock) GetManager() Manager {
	return m.manager
}

func (m *migrationMock) SetManager(manager Manager) Migration {
	m.manager = manager
	return m
}

func TestSourceCode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Source Class")
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
			d := NewCodeSource()
			d.Register(m1)
			d.Register(m2)

			Expect(d.migrations).To(HaveLen(2))
		})

		It("It should register migrations with inverted order", func() {
			d := NewCodeSource()
			d.Register(m1)
			d.Register(m2)
			d.Register(m3)

			Expect(d.migrations).To(HaveLen(3))
			Expect(d.migrations[0].GetDescription()).To(Equal("GetDescription 3"))
			Expect(d.migrations[1].GetDescription()).To(Equal("GetDescription 2"))
			Expect(d.migrations[2].GetDescription()).To(Equal("GetDescription 1"))
		})
	})
})
