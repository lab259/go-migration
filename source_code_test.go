package migration_test

import (
	"github.com/lab259/go-migration"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type migrationMock struct {
	id              time.Time
	description     string
	manager         migration.Manager
	done            bool
	undone          bool
	doneErr         error
	donePanicData   interface{}
	undoneErr       error
	undonePanicData interface{}
}

func (m *migrationMock) GetID() time.Time {
	return m.id
}

func (m *migrationMock) GetDescription() string {
	return m.description
}

func (m *migrationMock) Do() error {
	m.done = true
	if m.donePanicData != nil {
		panic(m.donePanicData)
	}
	return m.doneErr
}

func (m *migrationMock) Undo() error {
	m.undone = true
	if m.undonePanicData != nil {
		panic(m.undonePanicData)
	}
	return m.undoneErr
}

func (m *migrationMock) GetManager() migration.Manager {
	return m.manager
}

func (m *migrationMock) SetManager(manager migration.Manager) migration.Migration {
	m.manager = manager
	return m
}

var _ = Describe("Source Code", func() {
	var (
		m1 *migrationMock
		m2 *migrationMock
		m3 *migrationMock
	)

	BeforeEach(func() {
		m1 = &migrationMock{
			id:          time.Date(2001, 1, 1, 1, 1, 1, 1, time.UTC),
			description: "GetDescription 1",
		}

		m2 = &migrationMock{
			id:          time.Date(2001, 1, 1, 1, 1, 1, 0, time.UTC),
			description: "GetDescription 2",
		}

		m3 = &migrationMock{
			id:          time.Date(2000, 1, 1, 1, 1, 1, 0, time.UTC),
			description: "GetDescription 3",
		}
	})

	Describe("Register", func() {
		It("should register migrations", func() {
			d := migration.NewCodeSource()
			d.Register(m1)
			d.Register(m2)

			list, err := d.List()
			Expect(err).To(BeNil())
			Expect(list).To(HaveLen(2))
		})

		It("should register migrations ordered by ID", func() {
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

	Describe("List", func() {
		It("should list migrations", func() {
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
