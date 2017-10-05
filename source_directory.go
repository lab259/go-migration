package migration

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"regexp"
	"time"
	"errors"
	"fmt"
)

type DirectorySource struct {
	Directory string
	Extension string
}

var DirectorySourcePattern *regexp.Regexp = regexp.MustCompile("^([0-9]{14})_(.*)$")

func (this *DirectorySource) List() ([]Migration, error) {
	migrationsMap := make(map[string]*MigrationFile)
	if files, err := ioutil.ReadDir(this.Directory); err == nil {
		result := make([]Migration, 0)
		for i := 0; i < len(files); i++ {
			f := files[i]
			toks := strings.Split(filepath.Base(f.Name()), ".")
			if (len(toks) == 3) && (strings.ToLower(toks[len(toks)-1]) == strings.ToLower(this.Extension)) {
				var (
					id          time.Time
					description string
				)
				if tmpdata := DirectorySourcePattern.FindStringSubmatch(toks[0]); (len(tmpdata) != 3) || (tmpdata[1] == "") || (tmpdata[2] == "") {
					return nil, errors.New(fmt.Sprintf("%s does not meet the naming convention.", f.Name()))
				} else {
					id, err = time.Parse("20060102150405", tmpdata[1])
					description = strings.Replace(tmpdata[2], "_", " ", 0)
					if err != nil {
						return nil, err
					}
				}

				migration, ok := migrationsMap[toks[0]]
				if !ok {
					migration := &MigrationFile{
						id:          id,
						description: description,
						baseFile:    toks[0],
						up:          toks[1] == "up",
						down:        toks[1] == "down",
					}
					migrationsMap[migration.baseFile] = migration
					result = append(result, migration)
				} else if toks[1] == "up" {
					migration.up = true
				} else if toks[1] == "down" {
					migration.down = true
				}
			}
		}
		return result, nil
	} else {
		return nil, err
	}
}
