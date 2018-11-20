package migration_test

import (
	"github.com/globalsign/mgo"
	"github.com/lab259/go-migration"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

type migrationFromDB struct {
	ID time.Time `bson:"_id"`
}

func getSession() (*mgo.Session, error) {
	session, err := mgo.Dial("mongodb://localhost/test")
	if err != nil {
		return nil, err
	}
	return session, nil
}

var _ = Describe("MongoDBTarget", func() {
	var (
		session            *mgo.Session
		m1, m2, m3, m4, m5 *migration.DefaultMigration
		source             *migration.CodeSource
	)
	BeforeEach(func() {
		s, err := getSession()
		Expect(err).To(BeNil())
		session = s

		// Drops the _migrations table
		names, err := session.DB("").CollectionNames()
		Expect(err).ToNot(HaveOccurred())
		for _, name := range names {
			if name == migration.DefaultMigrationTable {
				err = session.DB("").C(migration.DefaultMigrationTable).DropCollection()
				Expect(err).ToNot(HaveOccurred())
			}
		}

		source = migration.NewCodeSource()

		baseTime := time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC)

		m1 = migration.NewMigration(baseTime, "Migration 1")
		m2 = migration.NewMigration(baseTime.Add(time.Second), "Migration 2")
		m3 = migration.NewMigration(baseTime.Add(time.Hour), "Migration 3")
		m4 = migration.NewMigration(baseTime.Add(time.Hour*24), "Migration 4")
		m5 = migration.NewMigration(baseTime.Add(time.Hour*24*10), "Migration 5")
		source.Register(m1)
		source.Register(m2)
		source.Register(m3)
		source.Register(m4)
		source.Register(m5)
	})

	AfterEach(func() {
		session.Close()
		session = nil
	})

	It("should return a new instance of a MongoDBTarget", func() {
		target := migration.NewMongoDB(session)
		Expect(target).NotTo(BeNil())
	})

	It("should add migrations to the execution list", func() {
		target := migration.NewMongoDB(session)
		target.AddMigration(migration.NewSummary(m1))
		target.AddMigration(migration.NewSummary(m2))
		target.AddMigration(migration.NewSummary(m5))

		c := session.DB("").C(migration.DefaultMigrationTable)
		migrations := make([]migrationFromDB, 0)
		err := c.Find(nil).Sort("_id").All(&migrations)
		Expect(err).ToNot(HaveOccurred())
		Expect(migrations).To(HaveLen(3))
		Expect(migrations[0].ID).To(Equal(m1.GetID()))
		Expect(migrations[1].ID).To(Equal(m2.GetID()))
		Expect(migrations[2].ID).To(Equal(m5.GetID()))
	})

	It("should return NoVersion when there is no migrations ran", func() {
		target := migration.NewMongoDB(session)
		version, err := target.Version()
		Expect(err).ToNot(HaveOccurred())
		Expect(version).To(Equal(migration.NoVersion))
	})

	It("should return the current version", func() {
		target := migration.NewMongoDB(session)
		target.AddMigration(migration.NewSummary(m1))
		target.AddMigration(migration.NewSummary(m2))
		target.AddMigration(migration.NewSummary(m5))

		version, err := target.Version()
		Expect(err).ToNot(HaveOccurred())
		Expect(version).To(Equal(m5.GetID()))
	})

	It("should return the current version with an arbitrary addition of migrations", func() {
		target := migration.NewMongoDB(session)
		target.AddMigration(migration.NewSummary(m5))
		target.AddMigration(migration.NewSummary(m3))
		target.AddMigration(migration.NewSummary(m1))

		version, err := target.Version()
		Expect(err).ToNot(HaveOccurred())
		Expect(version).To(Equal(m5.GetID()))
	})

	It("should return the current version with an arbitrary addition of migrations", func() {
		target := migration.NewMongoDB(session)
		target.AddMigration(migration.NewSummary(m5))
		target.AddMigration(migration.NewSummary(m3))
		target.AddMigration(migration.NewSummary(m1))

		migrations, err := target.MigrationsExecuted()
		Expect(err).ToNot(HaveOccurred())
		Expect(migrations).To(HaveLen(3))
		Expect(migrations[0]).To(Equal(m1.GetID()))
		Expect(migrations[1]).To(Equal(m3.GetID()))
		Expect(migrations[2]).To(Equal(m5.GetID()))
	})

	It("should remove a migration from the database", func() {
		target := migration.NewMongoDB(session)
		target.AddMigration(migration.NewSummary(m5))
		target.AddMigration(migration.NewSummary(m3))
		target.AddMigration(migration.NewSummary(m1))

		err := target.RemoveMigration(migration.NewSummary(m3))
		Expect(err).NotTo(HaveOccurred())

		migrations, err := target.MigrationsExecuted()
		Expect(err).ToNot(HaveOccurred())
		Expect(migrations).To(HaveLen(2))
		Expect(migrations[0]).To(Equal(m1.GetID()))
		Expect(migrations[1]).To(Equal(m5.GetID()))
	})
})
