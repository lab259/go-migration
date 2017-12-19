package migration

import (
	"flag"
	"fmt"
	"os"
)

type ArgsRunner struct {
	executor Executor
	args     []string
}

func NewArgsRunner(reporter Reporter, manager Manager, args ...string) Runner {
	return &ArgsRunner{
		executor: Executor{
			reporter: reporter,
			manager:  manager,
		},
		args: args,
	}
}

func NewArgsRunnerCustom(executor Executor, args ...string) Runner {
	return &ArgsRunner{
		executor: executor,
		args:     args,
	}
}

func (runner *ArgsRunner) Run() {
	fs := flag.NewFlagSet("Migrations", flag.ExitOnError)
	if err := fs.Parse(runner.args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, s := range flag.Args() {
		switch s {
		case "pending":
			runner.executor.Pending()
		case "executed":
			runner.executor.Executed()
		case "migrate":
			runner.executor.Migrate()
		case "rewind":
			runner.executor.Rewind()
		case "do":
			runner.executor.Do()
		case "undo":
			runner.executor.Undo()
		case "reset":
			runner.executor.Reset()
		default:
			// TODO runner.reporter.CommandNotFound(s)
			fmt.Println(fmt.Sprintf("\"%s\" command is not defined", s))
			os.Exit(2)
		}
		break
	}
}
