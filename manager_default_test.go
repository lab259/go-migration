package migration_test

import (
	"time"

	"github.com/lab259/go-migration"

	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type nopTarget struct {
	executed []*migration.Summary
}

// Version returns the current version of the database.
func (target *nopTarget) Version() (time.Time, error) {
	max := migration.NoVersion
	for _, m := range target.executed {
		if m.Migration.GetID().After(max) {
			max = m.Migration.GetID()
		}
	}
	return max, nil
}

// AddMigration persists the version on the database (or any other mean
// necessary).
func (target *nopTarget) AddMigration(summary *migration.Summary) error {
	if target.executed == nil {
		target.executed = make([]*migration.Summary, 0)
	}
	target.executed = append(target.executed, summary)
	return nil
}

// AddMigration persists the version on the database (or any other mean
// necessary).
func (target *nopTarget) RemoveMigration(summary *migration.Summary) error {
	if target.executed == nil {
		target.executed = make([]*migration.Summary, 0)
	}
	for i := len(target.executed) - 1; i >= 0; i-- {
		if target.executed[i].Migration.GetID() == summary.Migration.GetID() {
			target.executed = append(target.executed[:i], target.executed[i+1:]...)
		}
	}
	return nil
}

// MigrationsExecuted returns the IDs of the migrations executed
func (target *nopTarget) MigrationsExecuted() ([]time.Time, error) {
	executed := make([]time.Time, len(target.executed))
	for i, summary := range target.executed {
		executed[i] = summary.Migration.GetID()
	}
	return executed, nil
}

type ErroredTarget struct {
}

// Version returns NoVersion
func (target *ErroredTarget) Version() (time.Time, error) {
	return migration.NoVersion, nil
}

// AddMigration does nothing
func (target *ErroredTarget) AddMigration(summary *migration.Summary) error {
	return errors.New("AddMigration: forced error")
}

// RemoveMigration does nothing
func (target *ErroredTarget) RemoveMigration(summary *migration.Summary) error {
	return errors.New("RemoveMigration: forced error")
}

// MigrationsExecuted returns an error
func (target *ErroredTarget) MigrationsExecuted() ([]time.Time, error) {
	return nil, errors.New("MigrationsExecuted: forced error")
}

type AddMigrationErroredTarget struct {
	nopTarget
}

// AddMigration return an error
func (target *AddMigrationErroredTarget) AddMigration(summary *migration.Summary) error {
	return errors.New("AddMigration: forced error")
}

type RemoveMigrationErroredTarget struct {
	nopTarget
}

// RemoveMigration does nothing
func (target *RemoveMigrationErroredTarget) RemoveMigration(summary *migration.Summary) error {
	return errors.New("RemoveMigration: forced error")
}

type BeforeRunTarget struct {
	nopTarget
	BeforeRuns int
}

func (target *BeforeRunTarget) BeforeRun(executionContext interface{}) {
	target.BeforeRuns += 1
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

func (reporter *nopReporter) BeforeReset() {
}

func (reporter *nopReporter) AfterReset(rewindSummary []*migration.Summary, migrateSummary []*migration.Summary, err error) {
}

func (reporter *nopReporter) ListPending(migrations []migration.Migration, err error) {
}

func (reporter *nopReporter) ListExecuted(migrations []migration.Migration, err error) {
}

func (reporter *nopReporter) MigrationsStarved(migrations []migration.Migration) {
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
		/*m4UndoneErr, */ m5DoneErr *migrationMock
		/*m4UndonePanic,  m5DonePanic*/
	)

	BeforeEach(func() {
		target = &nopTarget{}

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

		/*
			m4UndoneErr = &migrationMock{
				id:          time.Date(2001, 6, 0, 0, 0, 0, 0, time.UTC),
				description: "GetDescription 4: Undone err",
				undoneErr:   errors.New("m4 undone forced error"),
			}
		*/

		m5DoneErr = &migrationMock{
			id:          time.Date(2001, 6, 0, 1, 0, 0, 0, time.UTC),
			description: "GetDescription 5: Done Err",
			doneErr:     errors.New("m5 done forced error"),
		}

		/*
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
		*/

		manager = migration.NewDefaultManager(target, codeSource)
	})

	It("should create a new instance of the DefaultManager", func() {
		t := &nopTarget{}
		cs := migration.NewCodeSource()
		manager := migration.NewDefaultManager(t, cs)
		Expect(manager).NotTo(BeNil())
		Expect(manager.Target()).To(Equal(t))
		Expect(manager.Source()).To(Equal(cs))
	})

	Describe("MigrationsPending", func() {
		It("should list pendent migrations when all migrations are new", func() {
			migrations, err := manager.MigrationsPending()
			Expect(err).To(BeNil())
			Expect(migrations).To(HaveLen(3))
		})

		It("should list pendent migrations when there is new and pendent migrations", func() {
			target.AddMigration(migration.NewSummary(m1))
			target.AddMigration(migration.NewSummary(m2))

			migrations, err := manager.MigrationsPending()
			Expect(err).To(BeNil())
			Expect(migrations).To(HaveLen(1))
			Expect(migrations[0].GetID()).To(Equal(m3.GetID()))
		})

		It("should list pendent migrations when there is new and pendent migrations 2", func() {
			target.AddMigration(migration.NewSummary(m1))

			migrations, err := manager.MigrationsPending()
			Expect(err).To(BeNil())
			Expect(migrations).To(HaveLen(2))
			Expect(migrations[0].GetID()).To(Equal(m2.GetID()))
			Expect(migrations[1].GetID()).To(Equal(m3.GetID()))
		})

		It("should list no pendents", func() {
			target.AddMigration(migration.NewSummary(m1))
			target.AddMigration(migration.NewSummary(m2))
			target.AddMigration(migration.NewSummary(m3))

			migrations, err := manager.MigrationsPending()
			Expect(err).To(BeNil())
			Expect(migrations).To(HaveLen(0))
		})
	})

	Describe("MigrationsExecuted", func() {
		It("should fail listing the migrations to be executed", func() {
			manager := migration.NewDefaultManager(&ErroredTarget{}, codeSource)
			migrations, err := manager.MigrationsExecuted()
			Expect(err).To(HaveOccurred())
			Expect(migrations).To(BeNil())
			Expect(err.Error()).To(Equal("MigrationsExecuted: forced error"))
		})

		It("should list executed migrations", func() {
			target.AddMigration(migration.NewSummary(m1))
			target.AddMigration(migration.NewSummary(m2))
			target.AddMigration(migration.NewSummary(m3))

			migrations, err := manager.MigrationsExecuted()
			Expect(err).To(BeNil())
			Expect(migrations).To(HaveLen(3))
			Expect(migrations[0].GetID()).To(Equal(m1.GetID()))
			Expect(migrations[1].GetID()).To(Equal(m2.GetID()))
			Expect(migrations[2].GetID()).To(Equal(m3.GetID()))
		})

		It("should list executed migrations migrations added in arbitrary order", func() {
			target.AddMigration(migration.NewSummary(m3))
			target.AddMigration(migration.NewSummary(m2))
			target.AddMigration(migration.NewSummary(m1))

			migrations, err := manager.MigrationsExecuted()
			Expect(err).To(BeNil())
			Expect(migrations).To(HaveLen(3))
			Expect(migrations[0].GetID()).To(Equal(m1.GetID()))
			Expect(migrations[1].GetID()).To(Equal(m2.GetID()))
			Expect(migrations[2].GetID()).To(Equal(m3.GetID()))
		})

		It("should list executed migrations when there is new and executed migrations", func() {
			target.AddMigration(migration.NewSummary(m1))

			migrations, err := manager.MigrationsExecuted()
			Expect(err).To(BeNil())
			Expect(migrations).To(HaveLen(1))
			Expect(migrations[0].GetID()).To(Equal(m1.GetID()))
		})

		It("should list executed migrations when there is new and executed migrations 2", func() {
			target.AddMigration(migration.NewSummary(m1))

			migrations, err := manager.MigrationsExecuted()
			Expect(err).To(BeNil())
			Expect(migrations).To(HaveLen(1))
			Expect(migrations[0].GetID()).To(Equal(m1.GetID()))
		})

		It("should list only the first one executed", func() {
			target.AddMigration(migration.NewSummary(m1))

			migrations, err := manager.MigrationsExecuted()
			Expect(err).To(BeNil())
			Expect(migrations).To(HaveLen(1))
			Expect(migrations[0].GetID()).To(Equal(m1.GetID()))
		})
	})

	Describe("Migrate", func() {
		It("should fail adding the migration as executed", func() {
			manager := migration.NewDefaultManager(&AddMigrationErroredTarget{}, codeSource)

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Migrate(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					migrations = append(migrations, summary.Migration)
				},
			}, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("AddMigration: forced error"))

			Expect(ms).To(HaveLen(1))
			Expect(ms[0].Migration).To(Equal(m1))

			Expect(migrations).To(HaveLen(1))
			Expect(migrations[0]).To(Equal(m1))

			Expect(m1.done).To(BeTrue())
		})

		It("should migrate all migrations", func() {
			migrations := make([]migration.Migration, 0)
			ms, err := manager.Migrate(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					migrations = append(migrations, summary.Migration)
				},
			}, nil)
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
			target.AddMigration(migration.NewSummary(m1))

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Migrate(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					migrations = append(migrations, summary.Migration)
				},
			}, nil)
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
			target.AddMigration(migration.NewSummary(m1))

			erroredMigration := &migrationMock{
				id:          m2.id.Add(time.Second),
				description: "Errored migration",
				doneErr:     errors.New("forced error"),
			}

			codeSource.Register(erroredMigration)

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Migrate(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					migrations = append(migrations, summary.Migration)
				},
			}, nil)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("forced error"))

			Expect(ms).To(HaveLen(2))
			Expect(ms[0].Migration.GetID()).To(Equal(m2.GetID()))
			Expect(ms[1].Migration.GetID()).To(Equal(erroredMigration.GetID()))
			Expect(ms[1].Failed()).To(BeTrue())
			Expect(ms[1].Failure()).ToNot(BeNil())
			Expect(ms[1].Failure().Error()).To(ContainSubstring("forced error"))

			Expect(migrations).To(HaveLen(2))
			Expect(migrations[0].GetID()).To(Equal(m2.GetID()))
			Expect(migrations[1].GetID()).To(Equal(erroredMigration.GetID()))

			Expect(m2.done).To(BeTrue())
			Expect(m2.undone).To(BeFalse())
			Expect(erroredMigration.done).To(BeTrue())
			Expect(erroredMigration.undone).To(BeFalse())

			Expect(manager.Target().Version()).To(Equal(m2.GetID()))
		})

		It("should panic the migration with a non error panic data", func() {
			now := time.Now().UTC()
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
			}, nil)
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

			Expect(manager.Target().Version()).To(Equal(migration.NoVersion))
		})

		It("should panic the migration with a panic data as an error", func() {
			target.AddMigration(migration.NewSummary(m1))
			codeSource = migration.NewCodeSource()
			migrationErrored := &migrationMock{
				id:            m1.GetID().Add(time.Second),
				donePanicData: errors.New("this is the panic data"),
			}
			codeSource.Register(migrationErrored)
			// target.AddMigration(migration.NewSummary(migrationErrored))
			manager := migration.NewDefaultManager(target, codeSource)

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Migrate(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					migrations = append(migrations, summary.Migration)
				},
			}, nil)
			Expect(err).To(Equal(migration.ErrMigrationPanicked))

			Expect(ms).To(HaveLen(1))
			Expect(ms[0].Migration).To(Equal(migrationErrored))
			Expect(ms[0].Failed()).To(BeTrue())
			Expect(ms[0].Failure()).ToNot(BeNil())
			Expect(ms[0].Failure().Error()).To(ContainSubstring("this is the panic data"))
			Expect(ms[0].Panicked()).To(BeTrue())
			Expect(ms[0].PanicData()).To(Equal(ms[0].Failure()))

			Expect(migrations).To(HaveLen(1))
			Expect(migrations[0]).To(Equal(migrationErrored))

			Expect(migrationErrored.done).To(BeTrue())
			Expect(migrationErrored.undone).To(BeFalse())

			Expect(manager.Target().Version()).To(Equal(m1.GetID()))
		})

		It("should detect starvation", func() {
			target.AddMigration(migration.NewSummary(m2))

			ms, err := manager.Migrate(&nopReporter{}, nil)
			Expect(err).To(Equal(migration.ErrMigrationStarved))

			Expect(ms).To(BeEmpty())

			Expect(manager.Target().Version()).To(Equal(m2.GetID()))
		})
	})

	Describe("Rewind", func() {
		It("should fail removing a migration from the list", func() {
			target := &RemoveMigrationErroredTarget{}
			manager := migration.NewDefaultManager(target, codeSource)
			target.AddMigration(migration.NewSummary(m1))
			target.AddMigration(migration.NewSummary(m2))
			target.AddMigration(migration.NewSummary(m3))

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Rewind(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					migrations = append(migrations, summary.Migration)
				},
			}, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("RemoveMigration: forced error"))

			Expect(ms).To(HaveLen(1))
			Expect(ms[0].Migration.GetID()).To(Equal(m3.GetID()))

			Expect(migrations).To(HaveLen(1))
			Expect(migrations[0].GetID()).To(Equal(m3.GetID()))

			Expect(m1.done).To(BeFalse())
			Expect(m1.undone).To(BeFalse())
			Expect(m2.done).To(BeFalse())
			Expect(m2.undone).To(BeFalse())
			Expect(m3.done).To(BeFalse())
			Expect(m3.undone).To(BeTrue())

			migrations, err = manager.MigrationsExecuted()
			Expect(err).ToNot(HaveOccurred())
			Expect(migrations).To(HaveLen(3))
			Expect(migrations[0].GetID()).To(Equal(m1.GetID()))
			Expect(migrations[1].GetID()).To(Equal(m2.GetID()))
			Expect(migrations[2].GetID()).To(Equal(m3.GetID()))
		})

		It("should rewind all migrations", func() {
			target.AddMigration(migration.NewSummary(m1))
			target.AddMigration(migration.NewSummary(m2))
			target.AddMigration(migration.NewSummary(m3))

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Rewind(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					migrations = append(migrations, summary.Migration)
				},
			}, nil)
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

			migrations, err = manager.MigrationsExecuted()
			Expect(err).ToNot(HaveOccurred())
			Expect(migrations).To(BeEmpty())
		})

		It("should rewind all migrations part of migrations", func() {
			target.AddMigration(migration.NewSummary(m1))
			target.AddMigration(migration.NewSummary(m2))

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Rewind(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					migrations = append(migrations, summary.Migration)
				},
			}, nil)
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
			target.AddMigration(migration.NewSummary(m1))
			target.AddMigration(migration.NewSummary(m2))
			target.AddMigration(migration.NewSummary(m3))
			migrationErrored := &migrationMock{
				id:          m2.id.Add(time.Second),
				description: "undone err",
				undoneErr:   errors.New("forced error"),
			}
			codeSource.Register(migrationErrored)
			target.AddMigration(migration.NewSummary(migrationErrored))

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Rewind(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					migrations = append(migrations, summary.Migration)
				},
			}, nil)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("forced error"))

			Expect(ms).To(HaveLen(2))
			Expect(ms[0].Migration).To(Equal(m3))
			Expect(ms[1].Migration).To(Equal(migrationErrored))
			Expect(ms[1].Failed()).To(BeTrue())
			Expect(ms[1].Failure()).ToNot(BeNil())
			Expect(ms[1].Failure().Error()).To(ContainSubstring("forced error"))

			Expect(migrations).To(HaveLen(2))
			Expect(migrations[0]).To(Equal(m3))
			Expect(migrations[1]).To(Equal(migrationErrored))

			Expect(migrationErrored.done).To(BeFalse())
			Expect(migrationErrored.undone).To(BeTrue())
			Expect(m3.done).To(BeFalse())
			Expect(m3.undone).To(BeTrue())

			Expect(manager.Target().Version()).To(Equal(migrationErrored.GetID()))
		})

		It("should not save the version when panicking", func() {
			target.AddMigration(migration.NewSummary(m1))
			target.AddMigration(migration.NewSummary(m2))
			target.AddMigration(migration.NewSummary(m3))
			migrationErrored := &migrationMock{
				id:              m2.id.Add(time.Second),
				description:     "GetDescription 4: Done Panic",
				undonePanicData: errors.New("m4 undone panic forced error"),
			}
			codeSource.Register(migrationErrored)
			target.AddMigration(migration.NewSummary(migrationErrored))

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Rewind(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					migrations = append(migrations, summary.Migration)
				},
			}, nil)
			Expect(err).To(Equal(migration.ErrMigrationPanicked))

			Expect(ms).To(HaveLen(2))
			Expect(ms[0].Migration).To(Equal(m3))
			Expect(ms[1].Migration).To(Equal(migrationErrored))
			Expect(ms[1].Panicked()).To(BeTrue())
			Expect(ms[1].Failed()).To(BeTrue())
			Expect(ms[1].Failure()).ToNot(BeNil())
			Expect(ms[1].Failure().Error()).To(ContainSubstring("m4 undone panic forced error"))

			Expect(migrations).To(HaveLen(2))
			Expect(migrations[0]).To(Equal(m3))
			Expect(migrations[1]).To(Equal(migrationErrored))

			Expect(m1.done).To(BeFalse())
			Expect(m1.undone).To(BeFalse())
			Expect(m2.done).To(BeFalse())
			Expect(m2.undone).To(BeFalse())
			Expect(migrationErrored.done).To(BeFalse())
			Expect(migrationErrored.undone).To(BeTrue())
			Expect(m3.done).To(BeFalse())
			Expect(m3.undone).To(BeTrue())

			Expect(manager.Target().Version()).To(Equal(migrationErrored.GetID()))
		})

		It("should panic the migration with a non error panic data", func() {
			target.AddMigration(migration.NewSummary(m1))
			codeSource = migration.NewCodeSource()
			migrationErrored := &migrationMock{
				id:              m1.id.Add(time.Second),
				undonePanicData: "this is the panic data",
			}
			codeSource.Register(migrationErrored)
			target.AddMigration(migration.NewSummary(migrationErrored))
			manager := migration.NewDefaultManager(target, codeSource)

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Rewind(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					migrations = append(migrations, summary.Migration)
				},
			}, nil)
			Expect(err).To(Equal(migration.ErrMigrationPanicked))

			Expect(ms).To(HaveLen(1))
			Expect(ms[0].Migration).To(Equal(migrationErrored))
			Expect(ms[0].Failed()).To(BeFalse())
			Expect(ms[0].Failure()).To(BeNil())
			Expect(ms[0].Panicked()).To(BeTrue())
			Expect(ms[0].PanicData()).To(Equal("this is the panic data"))

			Expect(migrations).To(HaveLen(1))
			Expect(migrations[0]).To(Equal(migrationErrored))

			Expect(migrationErrored.done).To(BeFalse())
			Expect(migrationErrored.undone).To(BeTrue())

			Expect(manager.Target().Version()).To(Equal(migrationErrored.GetID()))
		})

		It("should panic the migration with a panic data as an error", func() {
			target.AddMigration(migration.NewSummary(m1))
			codeSource = migration.NewCodeSource()
			migrationErrored := &migrationMock{
				id:              m1.GetID().Add(time.Second),
				undonePanicData: errors.New("this is the panic data"),
			}
			codeSource.Register(migrationErrored)
			target.AddMigration(migration.NewSummary(migrationErrored))
			manager := migration.NewDefaultManager(target, codeSource)

			migrations := make([]migration.Migration, 0)
			ms, err := manager.Rewind(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					migrations = append(migrations, summary.Migration)
				},
			}, nil)
			Expect(err).To(Equal(migration.ErrMigrationPanicked))

			Expect(ms).To(HaveLen(1))
			Expect(ms[0].Migration).To(Equal(migrationErrored))
			Expect(ms[0].Failed()).To(BeTrue())
			Expect(ms[0].Failure()).ToNot(BeNil())
			Expect(ms[0].Failure().Error()).To(ContainSubstring("this is the panic data"))
			Expect(ms[0].Panicked()).To(BeTrue())
			Expect(ms[0].PanicData()).To(Equal(ms[0].Failure()))

			Expect(migrations).To(HaveLen(1))
			Expect(migrations[0]).To(Equal(migrationErrored))

			Expect(migrationErrored.done).To(BeFalse())
			Expect(migrationErrored.undone).To(BeTrue())

			Expect(manager.Target().Version()).To(Equal(migrationErrored.GetID()))
		})
	})

	Describe("Do", func() {
		It("should execute a migration", func() {
			beforeMigrationCalled := false
			summary, err := manager.Do(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					Expect(summary.Migration).To(Equal(m1))
					beforeMigrationCalled = true
				},
			}, nil)

			Expect(err).To(BeNil())
			Expect(beforeMigrationCalled).To(BeTrue())
			Expect(summary.Direction()).To(Equal(migration.DirectionDo))
			Expect(summary.Migration).To(Equal(m1))

			version, err := target.Version()
			Expect(err).To(BeNil())
			Expect(version).To(Equal(m1.GetID()))

			migrationsExecuted, err := manager.MigrationsExecuted()
			Expect(err).ToNot(HaveOccurred())
			Expect(migrationsExecuted).To(HaveLen(1))
			Expect(migrationsExecuted[0].GetID()).To(Equal(m1.GetID()))
		})

		It("should multiple Do in sequence", func() {
			beforeMigrationCalled := false
			summary, err := manager.Do(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					Expect(summary.Migration).To(Equal(m1))
					beforeMigrationCalled = true
				},
			}, nil)

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
			}, nil)

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
			}, nil)

			Expect(err).To(BeNil())
			Expect(beforeMigrationCalled).To(BeTrue())
			Expect(summary.Direction()).To(Equal(migration.DirectionDo))
			Expect(summary.Migration).To(Equal(m3))

			version, err = target.Version()
			Expect(err).To(BeNil())
			Expect(version).To(Equal(m3.GetID()))

			migrationsExecuted, err := manager.MigrationsExecuted()
			Expect(err).ToNot(HaveOccurred())
			Expect(migrationsExecuted).To(HaveLen(3))
			Expect(migrationsExecuted[0].GetID()).To(Equal(m1.GetID()))
			Expect(migrationsExecuted[1].GetID()).To(Equal(m2.GetID()))
			Expect(migrationsExecuted[2].GetID()).To(Equal(m3.GetID()))
		})

		It("should fail migrating", func() {
			codeSource = migration.NewCodeSource()
			codeSource.Register(m5DoneErr)

			manager := migration.NewDefaultManager(target, codeSource)

			beforeMigrationCalled := false
			summary, err := manager.Do(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					Expect(summary.Migration).To(Equal(m5DoneErr))
					beforeMigrationCalled = true
				},
			}, nil)
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

			beforeMigrationCalled := false
			summary, err := manager.Do(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					Expect(summary.Migration).To(Equal(m))
					beforeMigrationCalled = true
				},
			}, nil)
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

			beforeMigrationCalled := false
			summary, err := manager.Do(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionDo))
					Expect(summary.Migration).To(Equal(m))
					beforeMigrationCalled = true
				},
			}, nil)
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

		It("should detect starvation", func() {
			target.AddMigration(migration.NewSummary(m2))

			ms, err := manager.Do(&nopReporter{}, nil)
			Expect(err).To(Equal(migration.ErrMigrationStarved))

			Expect(ms).To(BeNil())

			Expect(manager.Target().Version()).To(Equal(m2.GetID()))
		})
	})

	Describe("Undo", func() {
		It("should undo a migration", func() {
			manager.Target().AddMigration(migration.NewSummary(m1))

			beforeMigrationCalled := false
			summary, err := manager.Undo(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					Expect(summary.Migration).To(Equal(m1))
					beforeMigrationCalled = true
				},
			}, nil)

			Expect(err).To(BeNil())
			Expect(beforeMigrationCalled).To(BeTrue())
			Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
			Expect(summary.Migration).To(Equal(m1))

			version, err := target.Version()
			Expect(err).To(BeNil())
			Expect(version).To(Equal(migration.NoVersion))

			migrationsExecuted, err := manager.MigrationsExecuted()
			Expect(err).ToNot(HaveOccurred())
			Expect(migrationsExecuted).To(BeEmpty())
		})

		It("should multiple Undo in sequence", func() {
			manager.Target().AddMigration(migration.NewSummary(m1))
			manager.Target().AddMigration(migration.NewSummary(m2))
			manager.Target().AddMigration(migration.NewSummary(m3))

			beforeMigrationCalled := false
			summary, err := manager.Undo(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					Expect(summary.Migration).To(Equal(m3))
					beforeMigrationCalled = true
				},
			}, nil)

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
			}, nil)

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
			}, nil)

			Expect(err).To(BeNil())
			Expect(beforeMigrationCalled).To(BeTrue())
			Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
			Expect(summary.Migration).To(Equal(m1))

			version, err = target.Version()
			Expect(err).To(BeNil())
			Expect(version).To(Equal(migration.NoVersion))

			migrationsExecuted, err := manager.MigrationsExecuted()
			Expect(err).ToNot(HaveOccurred())
			Expect(migrationsExecuted).To(BeEmpty())
		})

		It("should fail migrating", func() {
			codeSource = migration.NewCodeSource()
			migrationErrored := &migrationMock{
				id:          m1.id.Add(time.Second),
				description: "Undone err",
				undoneErr:   errors.New("forced error"),
			}
			codeSource.Register(m1)
			codeSource.Register(migrationErrored)

			manager := migration.NewDefaultManager(target, codeSource)
			manager.Target().AddMigration(migration.NewSummary(m1))
			manager.Target().AddMigration(migration.NewSummary(migrationErrored))

			beforeMigrationCalled := false
			summary, err := manager.Undo(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					Expect(summary.Migration.GetID()).To(Equal(migrationErrored.id))
					beforeMigrationCalled = true
				},
			}, nil)
			Expect(err).To(Equal(migrationErrored.undoneErr))
			Expect(beforeMigrationCalled).To(BeTrue())
			Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
			Expect(summary.Migration).To(Equal(migrationErrored))
			Expect(summary.Failed()).To(BeTrue())
			Expect(summary.Failure()).To(Equal(migrationErrored.undoneErr))
			Expect(summary.Panicked()).To(BeFalse())
			Expect(summary.PanicData()).To(BeNil())

			version, err := target.Version()
			Expect(err).To(BeNil())
			Expect(version).To(Equal(migrationErrored.GetID()))

			migrationsExecuted, err := manager.MigrationsExecuted()
			Expect(migrationsExecuted).To(HaveLen(2))
			Expect(migrationsExecuted[0].GetID()).To(Equal(m1.GetID()))
			Expect(migrationsExecuted[1].GetID()).To(Equal(migrationErrored.GetID()))
		})

		It("should panic an error when migrating", func() {
			codeSource = migration.NewCodeSource()
			migrationErrored := &migrationMock{
				id:              m1.id.Add(time.Second),
				description:     "errored migration",
				undonePanicData: errors.New("forced error"),
			}
			codeSource.Register(m1)
			codeSource.Register(migrationErrored)

			manager := migration.NewDefaultManager(target, codeSource)
			target.AddMigration(migration.NewSummary(m1))
			target.AddMigration(migration.NewSummary(migrationErrored))

			beforeMigrationCalled := false
			summary, err := manager.Undo(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					Expect(summary.Migration).To(Equal(migrationErrored))
					beforeMigrationCalled = true
				},
			}, nil)
			Expect(beforeMigrationCalled).To(BeTrue())
			Expect(err).To(Equal(migration.ErrMigrationPanicked))
			Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
			Expect(summary.Migration).To(Equal(migrationErrored))
			Expect(summary.Failed()).To(BeTrue())
			Expect(summary.Failure()).ToNot(BeNil())
			Expect(summary.Failure().(error).Error()).To(ContainSubstring("forced error"))
			Expect(summary.Panicked()).To(BeTrue())
			Expect(summary.PanicData()).To(Equal(summary.Failure()))

			version, err := target.Version()
			Expect(err).To(BeNil())
			Expect(version).To(Equal(migrationErrored.GetID()))

			migrationsExecuted, err := manager.MigrationsExecuted()
			Expect(err).ToNot(HaveOccurred())
			Expect(migrationsExecuted).To(HaveLen(2))
			Expect(migrationsExecuted[0].GetID()).To(Equal(m1.GetID()))
			Expect(migrationsExecuted[1].GetID()).To(Equal(migrationErrored.GetID()))
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
			manager.Target().AddMigration(migration.NewSummary(m))

			beforeMigrationCalled := false
			summary, err := manager.Undo(&nopReporter{
				beforeMigration: func(summary *migration.Summary, err error) {
					Expect(summary.Direction()).To(Equal(migration.DirectionUndo))
					Expect(summary.Migration).To(Equal(m))
					beforeMigrationCalled = true
				},
			}, nil)
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

	Describe("Reset", func() {
		It("should reset a migration", func() {
			target.AddMigration(migration.NewSummary(m1))
			target.AddMigration(migration.NewSummary(m2))
			target.AddMigration(migration.NewSummary(m3))

			migrationUndone, migrationsDone, err := manager.Reset(&nopReporter{}, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(migrationUndone).To(HaveLen(3))
			Expect(migrationsDone).To(HaveLen(3))
		})
	})
})
