package migration

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// DirectorySource is migration.Source implementation. It provides the development
// of migrations using SQL files inside of a directory.
//
// The file pattern used to match the migrations is:
// <DATE>_<DESCRIPTION>.(up|down).(extension)
//
//     DATE: Must be in the format YYYYMMDDHHNNSS. For example, 20171005110647 for Oct 05, 2017 11:06:47)
//
//     DESCRIPTION: Any text but with no dots.
//
// Valid naming examples:
//
//     20171025191747_Creates_user_table.down.sql       : Drop the user table
//     20171025191747_Creates_user_table.up.sql         : Create user table
//     20171025191747_Creates customers table.down.sql  : Drop the customer table
//     20171025191747_Creates customers table.up.sql    : Create customer table
//     20171025213303_Drops_the_token_column.up.sql     : Drops a column, so is irreversible (there is no .down.sql)
//
// Invalid naming examples:
//     20171525191747_Creates_user_table.down.sql    : Invalid date (month 15?)
//     20171025191747_Creates_user_table.NNN.sql     : Invalid sufix (should be up or down)
type DirectorySource struct {
	// Directory represents the path that the migrations file will be searched
	// for.
	Directory string

	// Extension represents the file extension of the files.
	Extension string
}

var directorySourcePattern = regexp.MustCompile("^([0-9]{14})_(.*)$")

// List implements the migration.Source.List by listing all the files inside the
// migration.DirectorySource.Directory with the naming convention using the
// migration.DirectorySource.Extension.
func (s *DirectorySource) List() ([]Migration, error) {
	migrationsMap := make(map[string]*FileMigration)
	files, err := ioutil.ReadDir(s.Directory)
	if err == nil {
		result := make([]Migration, 0)
		for i := 0; i < len(files); i++ {
			f := files[i]
			toks := strings.Split(filepath.Base(f.Name()), ".")
			if (len(toks) == 3) && (strings.ToLower(toks[len(toks)-1]) == strings.ToLower(s.Extension)) {
				var (
					id          time.Time
					description string
				)
				if tmpdata := directorySourcePattern.FindStringSubmatch(toks[0]); !((len(tmpdata) != 3) || (tmpdata[1] == "") || (tmpdata[2] == "")) {
					id, err = time.Parse("20060102150405", tmpdata[1])
					description = strings.Replace(tmpdata[2], "_", " ", 0)
					if err != nil {
						return nil, err
					}
				} else {
					return nil, fmt.Errorf("%s does not meet the naming convention", f.Name())
				}

				migration, ok := migrationsMap[toks[0]]
				if !ok {
					migration := &FileMigration{
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
	}
	return nil, err
}
