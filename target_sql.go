package migration

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// MySQLTarget implements the migration.Target of the SQL based databases, using
// the Golang SQL package.
type MySQLTarget struct {
	connection *sql.DB
	tableName  string
}

// NewMySQL returns a new instance of the migration.MySQLTarget
func NewMySQL(conn *sql.DB) *MySQLTarget {
	return &MySQLTarget{
		tableName:  DefaultMigrationTable,
		connection: conn,
	}
}

// Version implements the migration.Target.Version by fetching the current
// version of the database from the table defined by
// migration.MongoDBTarget.SetCollectionName.
//
// It returns the current version of the database.
//
// Any error returned by the driver, will be passed to the caller.
func (target *MySQLTarget) Version() (time.Time, error) {
	ctx := context.Background()
	conn, err := target.connection.Conn(ctx)
	if err != nil {
		return NoVersion, err
	}
	defer conn.Close()

	err = target.ensureMigrationsTable(conn)
	if err != nil {
		return NoVersion, err
	}

	rows, err := conn.QueryContext(ctx, fmt.Sprintf("SELECT id FROM %s ORDER BY id DESC LIMIT 1", target.tableName))
	if err != nil {
		return time.Time{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return NoVersion, nil
	}
	d := time.Time{}
	err = rows.Scan(&d)
	if err != nil {
		return NoVersion, err
	}
	return d, nil
}

// SetVersion implements the migration.Target.SetVersion by storing the passed
// version on the database.
//
// It returns any error returned from the database driver.
func (target *MySQLTarget) SetVersion(id time.Time) error {
	ctx := context.Background()
	conn, err := target.connection.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = target.ensureMigrationsTable(conn)
	if err != nil {
		return err
	}

	_, err = conn.ExecContext(ctx, fmt.Sprintf("INSERT INTO %s (id) VALUES (?)", target.tableName), id)
	if err != nil {
		return err
	}
	return nil
}

// SetTableName sets the name of the table used to store the current migrations
// version that were executed.
func (target *MySQLTarget) SetTableName(collection string) *MySQLTarget {
	target.tableName = collection
	return target
}

// Connection returns the mgo.Session reference of this target.
func (target *MySQLTarget) Connection() *sql.DB {
	return target.connection
}

func (target *MySQLTarget) ensureMigrationsTable(conn *sql.Conn) error {
	ctx := context.Background()
	// ADHOC ADVISED: Tried to use the params with ? but it did not worked.
	rows, err := conn.QueryContext(ctx, fmt.Sprintf("SHOW TABLES LIKE \"%s\"", target.tableName))
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return nil
	}

	_, err = conn.ExecContext(ctx, fmt.Sprintf("CREATE TABLE %s (id DATETIME PRIMARY KEY)", target.tableName))
	return err
}
