package migration_test

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/lab259/go-migration"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func getDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "dbname=postgres user=postgres password=postgres host=localhost sslmode=disable")
	if err != nil {
		return nil, err
	}
	return db, nil
}

var _ = Describe("PostgreSQLTarget", func() {
	var (
		db                 *sql.DB
		m1, m2, m3, m4, m5 *migration.DefaultMigration
		source             *migration.CodeSource
	)
	BeforeEach(func() {
		d, err := getDB()
		Expect(err).ToNot(HaveOccurred())
		db = d

		// Drops the _migrations table
		_, err = db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %s`, pq.QuoteIdentifier(migration.DefaultMigrationTable)))
		Expect(err).ToNot(HaveOccurred())

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
		db.Close()
		db = nil
	})

	It("should return a new instance of a MongoDBTarget", func() {
		target := migration.NewPostgreSQLTarget(db)
		Expect(target).NotTo(BeNil())
	})

	It("should add migrations to the execution list", func() {
		target := migration.NewPostgreSQLTarget(db)
		target.AddMigration(migration.NewSummary(m1))
		target.AddMigration(migration.NewSummary(m2))
		target.AddMigration(migration.NewSummary(m5))

		migrations := make([]*migrationFromDB, 0, 3)
		rows, err := db.Query(fmt.Sprintf(`SELECT id FROM %s ORDER BY id`, pq.QuoteIdentifier(migration.DefaultMigrationTable)))
		Expect(err).ToNot(HaveOccurred())
		defer rows.Close()
		for rows.Next() {
			var id time.Time
			Expect(rows.Scan(&id)).To(Succeed())
			migrations = append(migrations, &migrationFromDB{id})
		}
		Expect(rows.Err()).ToNot(HaveOccurred())
		Expect(migrations).To(HaveLen(3))
		Expect(migrations[0].ID).To(Equal(m1.GetID()))
		Expect(migrations[1].ID).To(Equal(m2.GetID()))
		Expect(migrations[2].ID).To(Equal(m5.GetID()))
	})

	It("should return NoVersion when there is no migrations ran", func() {
		target := migration.NewPostgreSQLTarget(db)
		version, err := target.Version()
		Expect(err).ToNot(HaveOccurred())
		Expect(version).To(Equal(migration.NoVersion))
	})

	It("should return the current version", func() {
		target := migration.NewPostgreSQLTarget(db)
		target.AddMigration(migration.NewSummary(m1))
		target.AddMigration(migration.NewSummary(m2))
		target.AddMigration(migration.NewSummary(m5))

		version, err := target.Version()
		Expect(err).ToNot(HaveOccurred())
		Expect(version).To(Equal(m5.GetID()))
	})

	It("should return the current version with an arbitrary addition of migrations", func() {
		target := migration.NewPostgreSQLTarget(db)
		target.AddMigration(migration.NewSummary(m5))
		target.AddMigration(migration.NewSummary(m3))
		target.AddMigration(migration.NewSummary(m1))

		version, err := target.Version()
		Expect(err).ToNot(HaveOccurred())
		Expect(version).To(Equal(m5.GetID()))
	})

	It("should return the current version with an arbitrary addition of migrations", func() {
		target := migration.NewPostgreSQLTarget(db)
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
		target := migration.NewPostgreSQLTarget(db)
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
