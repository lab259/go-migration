package migration_test

import (
	"github.com/lab259/go-migration"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

type customReporter struct {
	beforeMigration  func(summary migration.Summary, err error)
	migrationSummary func(summary *migration.Summary, err error)
	afterMigration   func(summary migration.Summary, err error)
	beforeMigrate    func(migrations []migration.Migration)
	afterMigrate     func(migrations []*migration.Summary, err error)
	beforeRewind     func(migrations []migration.Migration)
	afterRewind      func(migrations []*migration.Summary, err error)
	beforeReset      func()
	afterReset       func(rewindSummary []*migration.Summary, migrateSummary []*migration.Summary, err error)
	listPending      func(migrations []migration.Migration, err error)
	listExecuted     func(migrations []migration.Migration, err error)
	migrationStarved func(migrations []migration.Migration)
	failure          func(err error)
	exit             func(code int)
	usage            func()
	commandNotFound  func(command string)
	noCommand        func()
}

func (reporter *customReporter) BeforeMigration(summary migration.Summary, err error) {
	if reporter.beforeMigration != nil {
		reporter.beforeMigration(summary, err)
	}
}

func (reporter *customReporter) MigrationSummary(summary *migration.Summary, err error) {
	if reporter.migrationSummary != nil {
		reporter.migrationSummary(summary, err)
	}
}

func (reporter *customReporter) AfterMigration(summary migration.Summary, err error) {
	if reporter.afterMigration != nil {
		reporter.afterMigration(summary, err)
	}
}

func (reporter *customReporter) BeforeMigrate(migrations []migration.Migration) {
	if reporter.beforeMigrate != nil {
		reporter.beforeMigrate(migrations)
	}
}

func (reporter *customReporter) AfterMigrate(migrations []*migration.Summary, err error) {
	if reporter.afterMigrate != nil {
		reporter.afterMigrate(migrations, err)
	}
}

func (reporter *customReporter) BeforeRewind(migrations []migration.Migration) {
	if reporter.beforeRewind != nil {
		reporter.beforeRewind(migrations)
	}
}

func (reporter *customReporter) AfterRewind(migrations []*migration.Summary, err error) {
	if reporter.afterRewind != nil {
		reporter.afterRewind(migrations, err)
	}
}

func (reporter *customReporter) BeforeReset() {
	if reporter.beforeReset != nil {
		reporter.beforeReset()
	}
}

func (reporter *customReporter) AfterReset(rewindSummary []*migration.Summary, migrateSummary []*migration.Summary, err error) {
	if reporter.afterReset != nil {
		reporter.afterReset(rewindSummary, migrateSummary, err)
	}
}

func (reporter *customReporter) ListPending(migrations []migration.Migration, err error) {
	if reporter.listPending != nil {
		reporter.listPending(migrations, err)
	}
}

func (reporter *customReporter) ListExecuted(migrations []migration.Migration, err error) {
	if reporter.listExecuted != nil {
		reporter.listExecuted(migrations, err)
	}
}

func (reporter *customReporter) MigrationsStarved(migrations []migration.Migration) {
	if reporter.migrationStarved != nil {
		reporter.migrationStarved(migrations)
	}
}

func (reporter *customReporter) Failure(err error) {
	if reporter.failure != nil {
		reporter.failure(err)
	}
}

func (reporter *customReporter) Exit(code int) {
	if reporter.exit != nil {
		reporter.exit(code)
	}
}

func (reporter *customReporter) Usage() {
	if reporter.usage != nil {
		reporter.usage()
	}
}

func (reporter *customReporter) CommandNotFound(command string) {
	if reporter.commandNotFound != nil {
		reporter.commandNotFound(command)
	}
}

func (reporter *customReporter) NoCommand() {
	if reporter.noCommand != nil {
		reporter.noCommand()
	}
}

var _ = Describe("RunnerArgs", func() {
	It("should create a new instance", func() {
		manager := migration.NewDefaultManager(&nopTarget{}, migration.NewCodeSource())
		r := migration.NewArgsRunner(&nopReporter{}, manager, func(code int) {
		})
		Expect(r).NotTo(BeNil())
	})

	It("should create a new instance with custom args", func() {
		manager := migration.NewDefaultManager(&nopTarget{}, migration.NewCodeSource())
		r := migration.NewArgsRunnerCustom(&nopReporter{}, manager, func(code int) {
		}, "migrate")
		Expect(r).NotTo(BeNil())
	})

	It("should run the pending command", func() {
		ran := false
		manager := migration.NewDefaultManager(&nopTarget{}, migration.NewCodeSource())
		r := migration.NewArgsRunnerCustom(&customReporter{
			listPending: func(migrations []migration.Migration, err error) {
				ran = true
			},
		}, manager, func(code int) {}, "pending")
		Expect(r).NotTo(BeNil())
		r.Run()
		Expect(ran).To(BeTrue())
	})

	It("should run the executed command", func() {
		ran := false
		manager := migration.NewDefaultManager(&nopTarget{}, migration.NewCodeSource())
		r := migration.NewArgsRunnerCustom(&customReporter{
			listExecuted: func(migrations []migration.Migration, err error) {
				ran = true
			},
		}, manager, func(code int) {}, "executed")
		Expect(r).NotTo(BeNil())
		r.Run()
		Expect(ran).To(BeTrue())
	})

	It("should run the migrate command", func() {
		ran := false
		manager := migration.NewDefaultManager(&nopTarget{}, migration.NewCodeSource())
		r := migration.NewArgsRunnerCustom(&customReporter{
			beforeMigrate: func(migrations []migration.Migration) {
				ran = true
			},
		}, manager, func(code int) {}, "migrate")
		Expect(r).NotTo(BeNil())
		r.Run()
		Expect(ran).To(BeTrue())
	})

	It("should run the rewind command", func() {
		ran := false
		manager := migration.NewDefaultManager(&nopTarget{}, migration.NewCodeSource())
		r := migration.NewArgsRunnerCustom(&customReporter{
			beforeRewind: func(migrations []migration.Migration) {
				ran = true
			},
		}, manager, func(code int) {}, "rewind")
		Expect(r).NotTo(BeNil())
		r.Run()
		Expect(ran).To(BeTrue())
	})

	It("should run the do command", func() {
		ran := false
		source := migration.NewCodeSource()
		source.Register(migration.NewMigration(time.Now(), "Description 1"))
		manager := migration.NewDefaultManager(&nopTarget{}, source)
		r := migration.NewArgsRunnerCustom(&customReporter{
			beforeMigration: func(summary migration.Summary, err error) {
				ran = true
			},
		}, manager, func(code int) {}, "do")
		Expect(r).NotTo(BeNil())
		r.Run()
		Expect(ran).To(BeTrue())
	})

	It("should run the undo command", func() {
		ran := false
		m := migration.NewMigration(time.Now(), "Description 1")
		source := migration.NewCodeSource()
		source.Register(m)
		target := &nopTarget{}
		target.AddMigration(migration.NewSummary(m))
		manager := migration.NewDefaultManager(target, source)
		r := migration.NewArgsRunnerCustom(&customReporter{
			beforeMigration: func(summary migration.Summary, err error) {
				ran = true
			},
		}, manager, func(code int) {}, "undo")
		Expect(r).NotTo(BeNil())
		r.Run()
		Expect(ran).To(BeTrue())
	})

	It("should run the reset command", func() {
		ran := false
		m := migration.NewMigration(time.Now(), "Description 1")
		source := migration.NewCodeSource()
		source.Register(m)
		target := &nopTarget{}
		target.AddMigration(migration.NewSummary(m))
		manager := migration.NewDefaultManager(target, source)
		r := migration.NewArgsRunnerCustom(&customReporter{
			beforeReset: func() {
				ran = true
			},
		}, manager, func(code int) {}, "reset")
		Expect(r).NotTo(BeNil())
		r.Run()
		Expect(ran).To(BeTrue())
	})

	It("should run an unknown command", func() {
		ran := false
		m := migration.NewMigration(time.Now(), "Description 1")
		source := migration.NewCodeSource()
		source.Register(m)
		target := &nopTarget{}
		target.AddMigration(migration.NewSummary(m))
		manager := migration.NewDefaultManager(target, source)
		r := migration.NewArgsRunnerCustom(&customReporter{
			commandNotFound: func(command string) {
				ran = true
				Expect(command).To(Equal("unknowncmd"))
			},
		}, manager, func(code int) {}, "unknowncmd")
		Expect(r).NotTo(BeNil())
		r.Run()
		Expect(ran).To(BeTrue())
	})

	It("should run an unknown command", func() {
		ran := false
		m := migration.NewMigration(time.Now(), "Description 1")
		source := migration.NewCodeSource()
		source.Register(m)
		target := &nopTarget{}
		target.AddMigration(migration.NewSummary(m))
		manager := migration.NewDefaultManager(target, source)
		r := migration.NewArgsRunnerCustom(&customReporter{
			noCommand: func() {
				ran = true
			},
		}, manager, func(code int) {})
		Expect(r).NotTo(BeNil())
		r.Run()
		Expect(ran).To(BeTrue())
	})
})
