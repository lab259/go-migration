package main

import (
	"bytes"
	"fmt"
	"github.com/lab259/go-migration"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "create" {
			create(args[1:]...)
			return
		} else {
			fmt.Println(args[0], "command not found")
		}
	}
	usage()
	os.Exit(3)
}

func create(description ...string) {
	if len(description) == 0 {
		fmt.Println("the description is required")
		fmt.Println("")
		usage()
	}

	descriptionBuff := bytes.NewBuffer(nil)
	for i, d := range description {
		if i > 0 {
			descriptionBuff.WriteString("_")
		}
		descriptionBuff.WriteString(d)
	}
	descriptionBuff.WriteString(".go")

	now := time.Now().UTC()
	nowstr := now.Format(migration.CodeMigrationDateFormat)
	fname := fmt.Sprintf("%s_%s", nowstr, descriptionBuff.String())
	_, err := os.Stat(fname)
	if err == nil {
		fmt.Println("This migration already exists.")
		os.Exit(1)
	}

	template, err := ioutil.ReadFile(".migration_template.go")
	if os.IsNotExist(err) {
		template = []byte(fmt.Sprintf(`package %s

import (
	"github.com/globalsign/mgo"
	. "github.com/lab259/go-migration"
)

func init() {
	NewCodeMigration(
		func() error {
			// Your migration up here
			
			return nil
		}, func() error {

			// Your migration down here

			return nil
		},
	)
}
`, path.Base(path.Dir(fname))))
	}

	f, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(template)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("%s created.", fname))
}

func usage() {
	fmt.Println("migration [command] [...params]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  create [description]    Create a new migrations in the current directory. If there is a")
	fmt.Println("                          .migration_template.go file it is used as template.")
	fmt.Println("                          + description: is the description of the migration.")
	fmt.Println("")
	fmt.Println("  help                    Displays this message.")
}
