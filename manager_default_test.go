package migration_test

import (
	"github.com/lab259/go-migration"
	"time"

	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type nopTarget struct {
	version time.Time
}

// Version returns the current version of the database.
func (target *nopTarget) Version() (time.Time, error) {
	return target.version, nil
}

// SetVersion persists the version on the database (or any other mean
// necessary).
func (target *nopTarget) SetVersion(version time.Time) error {
	target.version = version
	return nil
}

type nopReporter struct {
	beforeMigration func(summary *migration.Summary, err error)
}

func (reporter *nopReporter) BeforeMigration(summary migration.Summary, err error) {
	if reporter.beforeMigration != nil {
		reporter.beforeMigration(&summary, err)
	}
}

func (reporter *nopReporter) MigrationSummary(summary *migration.Summary, err error) {
}

func (reporter *nopReporter) AfterMigration(summary migration.Summary, err error) {
}

func (reporter *nopReporter) BeforeMigrate(migrations []migration.Migration) {
}

func (reporter *nopReporter) AfterMigrate(migrations []*migration.Summary, err error) {
}

func (reporter *nopReporter) BeforeRewind(migrations []migration.Migration) {
}

func (reporter *nopReporter) AfterRewind(migrations []*migration.Summary, err error) {
}

func (reporter *nopReporter) BeforeReset(doMigrations []migration.Migration, undoMigrations []migration.Migration) {
}

func (reporter *nopReporter) AfterReset(rewindSummary []*migration.Summary, migrateSummary []*migration.Summary, err error) {
}

func (reporter *nopReporter) ListPending(migrations []migration.Migration, err error) {
}

func (reporter *nopReporter) ListExecuted(migrations []migration.Migration, err error) {
}

func (reporter *nopReporter) Failure(err error) {
}

func (reporter *nopReporter) Exit(code int) {
}

func (reporter *nopReporter) Usage() {
}

func (reporter *nopReporter) CommandNotFound(command string) {
}

func (reporter *nopReporter) NoCommand() {
}

var _ = Describe("ManagerDefault", func() {

	var (
		target     migration.Target
		codeSource *migration.CodeSource
		manager    migration.Manager
		m1, m2, m3,
		m4UndoneErr, m5DoneErr,
		m4UndonePanic, m5DonePanic *migrationMock
	)

	BeforeEach(func() {
		target = &nopTarget{
			version: time.Now().UTC(),
		}

		m1 = &migrationMock{
			id:          time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC),
			description: "GetDescription 1",
		}

		m2 = &migrationMock{
			id:          time.Date(2001, 0, 0, 0, 0, 0, 0, time.UTC),
			description: "GetDescription 2",
		}

		m3 = &migrationMock{
			id:          time.Date(2002, 0, 0, 0, 0, 0, 0, time.UTC),
			description: "GetDescription 3",
		}

		codeSource = migration.NewCodeSource()
		codeSource.Register(m1)
		codeSource.Register(m2)
		codeSource.Register(m3)

		m4UndoneErr = &migrationMock{
			id:          time.Date(2001, 6, 0, 0, 0, 0, 0, time.UTC),
			description: "GetDescription 4: Undone err",
			undoneErr:   errors.New("m4 undone forced error"),
		}

		m5DoneErr = &migrationMock{
			id:          time.Date(2001, 6, 0, 1, 0, 0, 0, time.UTC),
			description: "GetDescription 5: Done Err",
			doneErr:     errors.New("m5 done forced error"),
		}

		m4UndonePanic = &migrationMock{
			id:              time.Date(2001, 6, 0, 0, 0, 0, 0, time.UTC),
			description:     "GetDescription 4: Done Panic",
			undonePanicData: errors.New("m4 undone panic forced error"),
		}

		m5DonePanic = &migrationMock{
			id:            time.Date(2001, 6, 0, 1, 0, 0, 0, time.UTC),
			description:   "GetDescription 5: Done Panic",
			donePanicData: errors.New("m5 done panic forced error"),
		}

		manager = migration.NewDefaultManager(target, codeSource)
	})

	It("should create a new instance of the DefaultManager", func() {
		t := &nopTarget{
			version: time.Now().UTC(),
		}
		cs := migration.NewCodeSource()
		manager := migration.NewDefaultManager(t, cs)
		Expect(manager).NotTo(BeNil())
		Expect(manager.Target()).To(Equal(t))
		Expect(manager.Source()).To(Equal(cs))
	})

	Describe("MigrationsPending", func() {
		It("should list pendent migrations when all migrations are new", func() {
			target.SetVersion(m1.GetID().Add(-time.Nanosecond))

			migrations, err := manager.MigrationsPending()
			Expect(err).To(BeNil())
			Expect(migrations).To(HaveLen(3))
		})

		It("should list pendent migrations when there is new and pendent migrations", func() {
			target.SetVersion(time.Date(2001, 1, 1, 1, 1, 1, 0, time.UTC))

			migrations, err := manager.MigrationsPending()
			Expect(err).To(BeNil())
			Expect(migrations).To(HaveLen(1))
			Expect(migrations[0].GetID()).To(Equal(m3.GetID()))
		})

		It("should list pendent migrations when there is new and pendent migrations 2", func() {
			target.SetVersion(m2.GetID().Add(-time.Nanosecond))

			migrations, err := manager.MigrationsPending()
			Expect(err).To(BeNil())
			Expect(migrations).To(HaveLen(2))
		})

		It("should list no pendents", func() {
			target.SetVersion(m3.GetID())

			migrations, err := manager.MigrationsPending()
			Expect(err).To(BeNil())
			Expect(migrations).To(HaveLen(0))
		})
	})

	Describe("MigrationsExecuted", func() {
		It("should list executed migrations when all migrations are old", func() {
			target.SetVersion(m3.GetID())

			migrations, err := manager.MigrationsExecuted()
			Expect(err).To(BeNil())
			Expect(migrations).To(HaveLen(3))
		})

		It("should list executed migrations migrations", func() {
			target.SetVersion(m3.GetID())

			migrations, err := manager.MigrationsExecuted()
			Expect(err).To(BeNil())
			Expect(migrations).To(HaveLen(3))
		})

		It("should list executed migrations when there is new and executed migrations", func() {
			target.SetVersion(time.Date(2000, 1, 1, 1, 1, 1, 0, time.UTC))

			migrations, err := manager.MigrationsExecuted()
			Expect(err).To(BeNil())
			Expect(migrations).To(HaveLen(1))
			Expect(migrations[0].GetID()).To(Equal(m1.GetID()))
		})

		It("should list executed migrations when there is new and executed migrations 2", func() {
			target.SetVersion(time.Date(2000, 1, 1, 1, 1, 3, 0, time.UTC))

			migrations, err := manager.MigrationsExecuted()
			Expect(err).To(BeNil())
			Expect(migrations).To(HaveLen(1))
			Expect(migrations[0].GetID()).To(Equal(m1.GetID()))
		})

		It("should list only the first one executed", func() {
			target.SetVersion(m1.GetID())

			migrations, err := manager.MigrationsExecuted()
			Expect(err).To(BeNil())
			Expect(migrations).To(HaveLen(1))
			Expect(migrations[0].GetID()).To(Equal(m1.GetID()))
		})
	})

	Describe("Migrate", func() {
		It("should migrate all migrations", func() {
			target.SetVersion(migration.NoVersion)

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Migrate(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					migrations = append(migrations, summary.Migration)
				},
			})
			Expect(err).To(BeNil())

			Expect(ms).To(HaveLen(3))
			Expect(ms[0].Migration).To(Equal(m1))
			Expect(ms[1].Migration).To(Equal(m2))
			Expect(ms[2].Migration).To(Equal(m3))

			Expect(migrations).To(HaveLen(3))
			Expect(migrations[0]).To(Equal(m1))
			Expect(migrations[1]).To(Equal(m2))
			Expect(migrations[2]).To(Equal(m3))

			Expect(m1.done).To(BeTrue())
			Expect(m2.done).To(BeTrue())
			Expect(m3.done).To(BeTrue())
		})

		It("should migrate only part of migrations", func() {
			target.SetVersion(m1.GetID())

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Migrate(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					migrations = append(migrations, summary.Migration)
				},
			})
			Expect(err).To(BeNil())

			Expect(ms).To(HaveLen(2))
			Expect(ms[0].Migration).To(Equal(m2))
			Expect(ms[1].Migration).To(Equal(m3))

			Expect(migrations).To(HaveLen(2))
			Expect(migrations[0]).To(Equal(m2))
			Expect(migrations[1]).To(Equal(m3))

			Expect(m1.done).To(BeFalse())
			Expect(m1.undone).To(BeFalse())
			Expect(m2.done).To(BeTrue())
			Expect(m2.undone).To(BeFalse())
			Expect(m3.done).To(BeTrue())
			Expect(m3.undone).To(BeFalse())

			Expect(manager.Target().Version()).To(Equal(m3.GetID()))
		})

		It("should migration fail in the middle", func() {
			target.SetVersion(m1.GetID().Add(-time.Hour))
			codeSource.Register(m4UndoneErr)
			codeSource.Register(m5DoneErr)

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Migrate(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					migrations = append(migrations, summary.Migration)
				},
			})
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("m5 done forced error"))

			Expect(ms).To(HaveLen(4))
			Expect(ms[0].Migration).To(Equal(m1))
			Expect(ms[1].Migration).To(Equal(m2))
			Expect(ms[2].Migration).To(Equal(m4UndoneErr))
			Expect(ms[3].Migration).To(Equal(m5DoneErr))
			Expect(ms[3].Failed()).To(BeTrue())
			Expect(ms[3].Failure()).ToNot(BeNil())
			Expect(ms[3].Failure().Error()).To(ContainSubstring("m5 done forced error"))

			Expect(migrations).To(HaveLen(4))
			Expect(migrations[0]).To(Equal(m1))
			Expect(migrations[1]).To(Equal(m2))
			Expect(migrations[2]).To(Equal(m4UndoneErr))
			Expect(migrations[3]).To(Equal(m5DoneErr))

			Expect(m1.done).To(BeTrue())
			Expect(m1.undone).To(BeFalse())
			Expect(m2.done).To(BeTrue())
			Expect(m2.undone).To(BeFalse())
			Expect(m4UndoneErr.done).To(BeTrue())
			Expect(m4UndoneErr.undone).To(BeFalse())
			Expect(m5DoneErr.done).To(BeTrue())
			Expect(m5DoneErr.undone).To(BeFalse())
			Expect(m3.done).To(BeFalse())
			Expect(m3.undone).To(BeFalse())

			Expect(manager.Target().Version()).To(Equal(m4UndoneErr.GetID()))
		})

		It("should panic the migration with a non error panic data", func() {
			now := time.Now().UTC()
			target.SetVersion(now.Add(-time.Nanosecond))
			codeSource = migration.NewCodeSource()
			m := &migrationMock{
				id:            now,
				donePanicData: "this is the panic data",
			}
			codeSource.Register(m)
			manager := migration.NewDefaultManager(target, codeSource)

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Migrate(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					migrations = append(migrations, summary.Migration)
				},
			})
			Expect(err).To(Equal(migration.ErrMigrationPanicked))

			Expect(ms).To(HaveLen(1))
			Expect(ms[0].Migration).To(Equal(m))
			Expect(ms[0].Failed()).To(BeFalse())
			Expect(ms[0].Failure()).To(BeNil())
			Expect(ms[0].Panicked()).To(BeTrue())
			Expect(ms[0].PanicData()).To(Equal("this is the panic data"))

			Expect(migrations).To(HaveLen(1))
			Expect(migrations[0]).To(Equal(m))

			Expect(m.done).To(BeTrue())
			Expect(m.undone).To(BeFalse())

			Expect(manager.Target().Version()).To(Equal(now.Add(-time.Nanosecond)))
		})

		It("should panic the migration with a panic data as an error", func() {
			now := time.Now().UTC()
			target.SetVersion(now.Add(-time.Nanosecond))
			codeSource = migration.NewCodeSource()
			m := &migrationMock{
				id:            now,
				donePanicData: errors.New("this is the panic data"),
			}
			codeSource.Register(m)
			manager := migration.NewDefaultManager(target, codeSource)

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Migrate(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					migrations = append(migrations, summary.Migration)
				},
			})
			Expect(err).To(Equal(migration.ErrMigrationPanicked))

			Expect(ms).To(HaveLen(1))
			Expect(ms[0].Migration).To(Equal(m))
			Expect(ms[0].Failed()).To(BeTrue())
			Expect(ms[0].Failure()).ToNot(BeNil())
			Expect(ms[0].Failure().Error()).To(ContainSubstring("this is the panic data"))
			Expect(ms[0].Panicked()).To(BeTrue())
			Expect(ms[0].PanicData()).To(Equal(ms[0].Failure()))

			Expect(migrations).To(HaveLen(1))
			Expect(migrations[0]).To(Equal(m))

			Expect(m.done).To(BeTrue())
			Expect(m.undone).To(BeFalse())

			Expect(manager.Target().Version()).To(Equal(now.Add(-time.Nanosecond)))
		})

		It("should not save version when panicking", func() {
			target.SetVersion(m1.GetID().Add(-time.Hour))
			codeSource.Register(m4UndonePanic)
			codeSource.Register(m5DonePanic)

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Migrate(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					migrations = append(migrations, summary.Migration)
				},
			})
			Expect(err).To(Equal(migration.ErrMigrationPanicked))

			Expect(ms).To(HaveLen(4))
			Expect(ms[0].Migration).To(Equal(m1))
			Expect(ms[1].Migration).To(Equal(m2))
			Expect(ms[2].Migration).To(Equal(m4UndonePanic))
			Expect(ms[3].Migration).To(Equal(m5DonePanic))
			Expect(ms[3].Failed()).To(BeTrue())
			Expect(ms[3].Failure()).ToNot(BeNil())
			Expect(ms[3].Failure().Error()).To(ContainSubstring("m5 done panic forced error"))

			Expect(migrations).To(HaveLen(4))
			Expect(migrations[0]).To(Equal(m1))
			Expect(migrations[1]).To(Equal(m2))
			Expect(migrations[2]).To(Equal(m4UndonePanic))
			Expect(migrations[3]).To(Equal(m5DonePanic))

			Expect(m1.done).To(BeTrue())
			Expect(m1.undone).To(BeFalse())
			Expect(m2.done).To(BeTrue())
			Expect(m2.undone).To(BeFalse())
			Expect(m4UndonePanic.done).To(BeTrue())
			Expect(m4UndonePanic.undone).To(BeFalse())
			Expect(m5DonePanic.done).To(BeTrue())
			Expect(m5DonePanic.undone).To(BeFalse())
			Expect(m3.done).To(BeFalse())
			Expect(m3.undone).To(BeFalse())

			Expect(manager.Target().Version()).To(Equal(m4UndoneErr.GetID()))
		})
	})

	Describe("Rewind", func() {
		It("should rewind all migrations", func() {
			target.SetVersion(m3.GetID())

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Rewind(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					migrations = append(migrations, summary.Migration)
				},
			})
			Expect(err).To(BeNil())

			Expect(ms).To(HaveLen(3))
			Expect(ms[0].Migration).To(Equal(m3))
			Expect(ms[1].Migration).To(Equal(m2))
			Expect(ms[2].Migration).To(Equal(m1))

			Expect(migrations).To(HaveLen(3))
			Expect(migrations[0]).To(Equal(m3))
			Expect(migrations[1]).To(Equal(m2))
			Expect(migrations[2]).To(Equal(m1))

			Expect(m1.done).To(BeFalse())
			Expect(m1.undone).To(BeTrue())
			Expect(m2.done).To(BeFalse())
			Expect(m2.undone).To(BeTrue())
			Expect(m3.done).To(BeFalse())
			Expect(m3.undone).To(BeTrue())
		})

		It("should rewind all migrations part of migrations", func() {
			target.SetVersion(m2.GetID())

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Rewind(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					migrations = append(migrations, summary.Migration)
				},
			})
			Expect(err).To(BeNil())

			Expect(ms).To(HaveLen(2))
			Expect(ms[0].Migration).To(Equal(m2))
			Expect(ms[1].Migration).To(Equal(m1))

			Expect(migrations).To(HaveLen(2))
			Expect(migrations[0]).To(Equal(m2))
			Expect(migrations[1]).To(Equal(m1))

			Expect(m1.done).To(BeFalse())
			Expect(m1.undone).To(BeTrue())
			Expect(m2.done).To(BeFalse())
			Expect(m2.undone).To(BeTrue())
			Expect(m3.done).To(BeFalse())
			Expect(m3.undone).To(BeFalse())
		})

		It("should rewind fail in the middle", func() {
			target.SetVersion(m3.GetID())
			codeSource.Register(m4UndoneErr)
			codeSource.Register(m5DoneErr)

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Rewind(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					migrations = append(migrations, summary.Migration)
				},
			})
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("m4 undone forced error"))

			Expect(ms).To(HaveLen(3))
			Expect(ms[0].Migration).To(Equal(m3))
			Expect(ms[1].Migration).To(Equal(m5DoneErr))
			Expect(ms[2].Migration).To(Equal(m4UndoneErr))
			Expect(ms[2].Failed()).To(BeTrue())
			Expect(ms[2].Failure()).ToNot(BeNil())
			Expect(ms[2].Failure().Error()).To(ContainSubstring("m4 undone forced error"))

			Expect(migrations).To(HaveLen(3))
			Expect(migrations[0]).To(Equal(m3))
			Expect(migrations[1]).To(Equal(m5DoneErr))
			Expect(migrations[2]).To(Equal(m4UndoneErr))

			Expect(m1.done).To(BeFalse())
			Expect(m1.undone).To(BeFalse())
			Expect(m2.done).To(BeFalse())
			Expect(m2.undone).To(BeFalse())
			Expect(m4UndoneErr.done).To(BeFalse())
			Expect(m4UndoneErr.undone).To(BeTrue())
			Expect(m5DoneErr.done).To(BeFalse())
			Expect(m5DoneErr.undone).To(BeTrue())
			Expect(m3.done).To(BeFalse())
			Expect(m3.undone).To(BeTrue())

			Expect(manager.Target().Version()).To(Equal(m5DoneErr.GetID()))
		})

		It("should not save the version when panicking", func() {
			target.SetVersion(m3.GetID())
			codeSource.Register(m4UndonePanic)
			codeSource.Register(m5DonePanic)

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Rewind(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					migrations = append(migrations, summary.Migration)
				},
			})
			Expect(err).To(Equal(migration.ErrMigrationPanicked))

			Expect(ms).To(HaveLen(3))
			Expect(ms[0].Migration).To(Equal(m3))
			Expect(ms[1].Migration).To(Equal(m5DonePanic))
			Expect(ms[2].Migration).To(Equal(m4UndonePanic))
			Expect(ms[2].Panicked()).To(BeTrue())
			Expect(ms[2].Failed()).To(BeTrue())
			Expect(ms[2].Failure()).ToNot(BeNil())
			Expect(ms[2].Failure().Error()).To(ContainSubstring("m4 undone panic forced error"))

			Expect(migrations).To(HaveLen(3))
			Expect(migrations[0]).To(Equal(m3))
			Expect(migrations[1]).To(Equal(m5DonePanic))
			Expect(migrations[2]).To(Equal(m4UndonePanic))

			Expect(m1.done).To(BeFalse())
			Expect(m1.undone).To(BeFalse())
			Expect(m2.done).To(BeFalse())
			Expect(m2.undone).To(BeFalse())
			Expect(m4UndonePanic.done).To(BeFalse())
			Expect(m4UndonePanic.undone).To(BeTrue())
			Expect(m5DonePanic.done).To(BeFalse())
			Expect(m5DonePanic.undone).To(BeTrue())
			Expect(m3.done).To(BeFalse())
			Expect(m3.undone).To(BeTrue())

			Expect(manager.Target().Version()).To(Equal(m5DoneErr.GetID()))
		})

		It("should panic the migration with a non error panic data", func() {
			now := time.Now().UTC()
			target.SetVersion(now.Add(time.Nanosecond))
			codeSource = migration.NewCodeSource()
			m := &migrationMock{
				id:              now,
				undonePanicData: "this is the panic data",
			}
			codeSource.Register(m)
			manager := migration.NewDefaultManager(target, codeSource)

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Rewind(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					migrations = append(migrations, summary.Migration)
				},
			})
			Expect(err).To(Equal(migration.ErrMigrationPanicked))

			Expect(ms).To(HaveLen(1))
			Expect(ms[0].Migration).To(Equal(m))
			Expect(ms[0].Failed()).To(BeFalse())
			Expect(ms[0].Failure()).To(BeNil())
			Expect(ms[0].Panicked()).To(BeTrue())
			Expect(ms[0].PanicData()).To(Equal("this is the panic data"))

			Expect(migrations).To(HaveLen(1))
			Expect(migrations[0]).To(Equal(m))

			Expect(m.done).To(BeFalse())
			Expect(m.undone).To(BeTrue())

			Expect(manager.Target().Version()).To(Equal(now.Add(time.Nanosecond)))
		})

		It("should panic the migration with a panic data as an error", func() {
			now := time.Now().UTC()
			target.SetVersion(now.Add(time.Nanosecond))
			codeSource = migration.NewCodeSource()
			m := &migrationMock{
				id:              now,
				undonePanicData: errors.New("this is the panic data"),
			}
			codeSource.Register(m)
			manager := migration.NewDefaultManager(target, codeSource)

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Rewind(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					migrations = append(migrations, summary.Migration)
				},
			})
			Expect(err).To(Equal(migration.ErrMigrationPanicked))

			Expect(ms).To(HaveLen(1))
			Expect(ms[0].Migration).To(Equal(m))
			Expect(ms[0].Failed()).To(BeTrue())
			Expect(ms[0].Failure()).ToNot(BeNil())
			Expect(ms[0].Failure().Error()).To(ContainSubstring("this is the panic data"))
			Expect(ms[0].Panicked()).To(BeTrue())
			Expect(ms[0].PanicData()).To(Equal(ms[0].Failure()))

			Expect(migrations).To(HaveLen(1))
			Expect(migrations[0]).To(Equal(m))

			Expect(m.done).To(BeFalse())
			Expect(m.undone).To(BeTrue())

			Expect(manager.Target().Version()).To(Equal(now.Add(time.Nanosecond)))
		})
	})

	Describe("Do", func() {
		It("should execute a migration", func() {
			manager.Target().SetVersion(migration.NoVersion)

			beforeMigrationCalled := false
			summary, err := manager.Do(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					Expect(summary.Migration).To(Equal(m1))
					beforeMigrationCalled = true
				},
			})

			Expect(err).To(BeNil())
			Expect(beforeMigrationCalled).To(BeTrue())
			Expect(summary.Direction()).To(Equal(migration.DirectionDo))
			Expect(summary.Migration).To(Equal(m1))

			version, err := target.Version()
			Expect(err).To(BeNil())
			Expect(version).To(Equal(m1.GetID()))
		})

		It("should multiple Do in sequence", func() {
			manager.Target().SetVersion(migration.NoVersion)

			beforeMigrationCalled := false
			summary, err := manager.Do(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					Expect(summary.Migration).To(Equal(m1))
					beforeMigrationCalled = true
				},
			})

			Expect(err).To(BeNil())
			Expect(beforeMigrationCalled).To(BeTrue())
			Expect(summary.Direction()).To(Equal(migration.DirectionDo))
			Expect(summary.Migration).To(Equal(m1))

			version, err := target.Version()
			Expect(err).To(BeNil())
			Expect(version).To(Equal(m1.GetID()))

			// Second
			beforeMigrationCalled = false
			summary, err = manager.Do(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					Expect(summary.Migration).To(Equal(m2))
					beforeMigrationCalled = true
				},
			})

			Expect(err).To(BeNil())
			Expect(beforeMigrationCalled).To(BeTrue())
			Expect(summary.Direction()).To(Equal(migration.DirectionDo))
			Expect(summary.Migration).To(Equal(m2))

			version, err = target.Version()
			Expect(err).To(BeNil())
			Expect(version).To(Equal(m2.GetID()))

			// Third
			beforeMigrationCalled = false
			summary, err = manager.Do(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					Expect(summary.Migration).To(Equal(m3))
					beforeMigrationCalled = true
				},
			})

			Expect(err).To(BeNil())
			Expect(beforeMigrationCalled).To(BeTrue())
			Expect(summary.Direction()).To(Equal(migration.DirectionDo))
			Expect(summary.Migration).To(Equal(m3))

			version, err = target.Version()
			Expect(err).To(BeNil())
			Expect(version).To(Equal(m3.GetID()))
		})

		It("should fail migrating", func() {
			codeSource = migration.NewCodeSource()
			codeSource.Register(m5DoneErr)

			manager := migration.NewDefaultManager(target, codeSource)
			manager.Target().SetVersion(migration.NoVersion)

			beforeMigrationCalled := false
			summary, err := manager.Do(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					Expect(summary.Migration).To(Equal(m5DoneErr))
					beforeMigrationCalled = true
				},
			})
			Expect(err).To(Equal(m5DoneErr.doneErr))
			Expect(beforeMigrationCalled).To(BeTrue())
			Expect(summary.Direction()).To(Equal(migration.DirectionDo))
			Expect(summary.Migration).To(Equal(m5DoneErr))
			Expect(summary.Failed()).To(BeTrue())
			Expect(summary.Failure()).To(Equal(m5DoneErr.doneErr))
			Expect(summary.Panicked()).To(BeFalse())
			Expect(summary.PanicData()).To(BeNil())

			version, err := target.Version()
			Expect(err).To(BeNil())
			Expect(version).To(Equal(migration.NoVersion))
		})

		It("should panic an error when migrating", func() {
			codeSource = migration.NewCodeSource()
			now := time.Now().UTC()
			m := &migrationMock{
				id:            now,
				donePanicData: errors.New("forced error"),
			}
			codeSource.Register(m)

			manager := migration.NewDefaultManager(target, codeSource)
			manager.Target().SetVersion(migration.NoVersion)

			beforeMigrationCalled := false
			summary, err := manager.Do(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					Expect(summary.Migration).To(Equal(m))
					beforeMigrationCalled = true
				},
			})
			Expect(err).To(Equal(migration.ErrMigrationPanicked))
			Expect(beforeMigrationCalled).To(BeTrue())
			Expect(summary.Direction()).To(Equal(migration.DirectionDo))
			Expect(summary.Migration).To(Equal(m))
			Expect(summary.Failed()).To(BeTrue())
			Expect(summary.Failure()).ToNot(BeNil())
			Expect(summary.Failure().(error).Error()).To(ContainSubstring("forced error"))
			Expect(summary.Panicked()).To(BeTrue())
			Expect(summary.PanicData()).To(Equal(summary.Failure()))

			version, err := target.Version()
			Expect(err).To(BeNil())
			Expect(version).To(Equal(migration.NoVersion))
		})

		It("should panic an string when migrating", func() {
			codeSource = migration.NewCodeSource()
			now := time.Now().UTC()
			m := &migrationMock{
				id:            now,
				donePanicData: "panicked data",
			}
			codeSource.Register(m)

			manager := migration.NewDefaultManager(target, codeSource)
			manager.Target().SetVersion(migration.NoVersion)

			beforeMigrationCalled := false
			summary, err := manager.Do(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					Expect(summary.Migration).To(Equal(m))
					beforeMigrationCalled = true
				},
			})
			Expect(err).To(Equal(migration.ErrMigrationPanicked))
			Expect(beforeMigrationCalled).To(BeTrue())
			Expect(summary.Direction()).To(Equal(migration.DirectionDo))
			Expect(summary.Migration).To(Equal(m))
			Expect(summary.Failed()).To(BeFalse())
			Expect(summary.Failure()).To(BeNil())
			Expect(summary.Panicked()).To(BeTrue())
			Expect(summary.PanicData()).To(Equal("panicked data"))

			version, err := target.Version()
			Expect(err).To(BeNil())
			Expect(version).To(Equal(migration.NoVersion))
		})
	})

	Describe("Undo", func() {
		It("should undo a migration", func() {
			manager.Target().SetVersion(m1.GetID())

			beforeMigrationCalled := false
			summary, err := manager.Undo(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					Expect(summary.Migration).To(Equal(m1))
					beforeMigrationCalled = true
				},
			})

			Expect(err).To(BeNil())
			Expect(beforeMigrationCalled).To(BeTrue())
			Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
			Expect(summary.Migration).To(Equal(m1))

			version, err := target.Version()
			Expect(err).To(BeNil())
			Expect(version).To(Equal(migration.NoVersion))
		})

		It("should multiple Undo in sequence", func() {
			manager.Target().SetVersion(m3.GetID())

			beforeMigrationCalled := false
			summary, err := manager.Undo(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					Expect(summary.Migration).To(Equal(m3))
					beforeMigrationCalled = true
				},
			})

			Expect(err).To(BeNil())
			Expect(beforeMigrationCalled).To(BeTrue())
			Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
			Expect(summary.Migration).To(Equal(m3))

			version, err := target.Version()
			Expect(err).To(BeNil())
			Expect(version).To(Equal(m2.GetID()))

			// Second
			beforeMigrationCalled = false
			summary, err = manager.Undo(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					Expect(summary.Migration).To(Equal(m2))
					beforeMigrationCalled = true
				},
			})

			Expect(err).To(BeNil())
			Expect(beforeMigrationCalled).To(BeTrue())
			Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
			Expect(summary.Migration).To(Equal(m2))

			version, err = target.Version()
			Expect(err).To(BeNil())
			Expect(version).To(Equal(m1.GetID()))

			// Third
			beforeMigrationCalled = false
			summary, err = manager.Undo(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					Expect(summary.Migration).To(Equal(m1))
					beforeMigrationCalled = true
				},
			})

			Expect(err).To(BeNil())
			Expect(beforeMigrationCalled).To(BeTrue())
			Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
			Expect(summary.Migration).To(Equal(m1))

			version, err = target.Version()
			Expect(err).To(BeNil())
			Expect(version).To(Equal(migration.NoVersion))
		})

		It("should fail migrating", func() {
			codeSource = migration.NewCodeSource()
			codeSource.Register(m4UndoneErr)

			manager := migration.NewDefaultManager(target, codeSource)
			manager.Target().SetVersion(m4UndoneErr.GetID())

			beforeMigrationCalled := false
			summary, err := manager.Undo(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					Expect(summary.Migration).To(Equal(m4UndoneErr))
					beforeMigrationCalled = true
				},
			})
			Expect(err).To(Equal(m4UndoneErr.undoneErr))
			Expect(beforeMigrationCalled).To(BeTrue())
			Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
			Expect(summary.Migration).To(Equal(m4UndoneErr))
			Expect(summary.Failed()).To(BeTrue())
			Expect(summary.Failure()).To(Equal(m4UndoneErr.undoneErr))
			Expect(summary.Panicked()).To(BeFalse())
			Expect(summary.PanicData()).To(BeNil())

			version, err := target.Version()
			Expect(err).To(BeNil())
			Expect(version).To(Equal(m4UndoneErr.GetID()))
		})

		It("should panic an error when migrating", func() {
			codeSource = migration.NewCodeSource()
			now := time.Now().UTC()
			m := &migrationMock{
				id:              now,
				undonePanicData: errors.New("forced error"),
			}
			codeSource.Register(m)

			manager := migration.NewDefaultManager(target, codeSource)
			manager.Target().SetVersion(now)

			beforeMigrationCalled := false
			summary, err := manager.Undo(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					Expect(summary.Migration).To(Equal(m))
					beforeMigrationCalled = true
				},
			})
			Expect(beforeMigrationCalled).To(BeTrue())
			Expect(err).To(Equal(migration.ErrMigrationPanicked))
			Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
			Expect(summary.Migration).To(Equal(m))
			Expect(summary.Failed()).To(BeTrue())
			Expect(summary.Failure()).ToNot(BeNil())
			Expect(summary.Failure().(error).Error()).To(ContainSubstring("forced error"))
			Expect(summary.Panicked()).To(BeTrue())
			Expect(summary.PanicData()).To(Equal(summary.Failure()))

			version, err := target.Version()
			Expect(err).To(BeNil())
			Expect(version).To(Equal(m.GetID()))
		})

		It("should panic an string when migrating", func() {
			codeSource = migration.NewCodeSource()
			now := time.Now().UTC()
			m := &migrationMock{
				id:              now,
				undonePanicData: "panicked data",
			}
			codeSource.Register(m)

			manager := migration.NewDefaultManager(target, codeSource)
			manager.Target().SetVersion(now)

			beforeMigrationCalled := false
			summary, err := manager.Undo(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					Expect(summary.Migration).To(Equal(m))
					beforeMigrationCalled = true
				},
			})
			Expect(err).To(Equal(migration.ErrMigrationPanicked))
			Expect(beforeMigrationCalled).To(BeTrue())
			Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
			Expect(summary.Migration).To(Equal(m))
			Expect(summary.Failed()).To(BeFalse())
			Expect(summary.Failure()).To(BeNil())
			Expect(summary.Panicked()).To(BeTrue())
			Expect(summary.PanicData()).To(Equal("panicked data"))

			version, err := target.Version()
			Expect(err).To(BeNil())
			Expect(version).To(Equal(m.GetID()))
		})
	})
})
