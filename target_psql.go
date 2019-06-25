package migration

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type PostgreSQLTarget struct {
	db        *sql.DB
	tableName string
}

func NewPostgreSQLTarget(db *sql.DB) *PostgreSQLTarget {
	return &PostgreSQLTarget{
		db:        db,
		tableName: pq.QuoteIdentifier(DefaultMigrationTable),
	}
}

func (target *PostgreSQLTarget) Version() (time.Time, error) {
	version := NoVersion
	err := target.withConn(func(ctx context.Context, conn *sql.Conn) error {
		rows, err := conn.QueryContext(ctx, fmt.Sprintf("SELECT id FROM %s ORDER BY id DESC LIMIT 1", target.tableName))
		if err != nil {
			return err
		}
		defer rows.Close()

		if !rows.Next() {
			return nil
		}

		if err := rows.Scan(&version); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return NoVersion, err
	}
	return version, nil
}

func (target *PostgreSQLTarget) AddMigration(summary *Summary) error {
	return target.withConn(func(ctx context.Context, conn *sql.Conn) error {
		_, err := conn.ExecContext(ctx, fmt.Sprintf("INSERT INTO %s (id) values ($1)", target.tableName), summary.Migration.GetID())
		return err
	})
}

func (target *PostgreSQLTarget) RemoveMigration(summary *Summary) error {
	return target.withConn(func(ctx context.Context, conn *sql.Conn) error {
		_, err := conn.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = $1", target.tableName), summary.Migration.GetID())
		return err
	})
}

func (target *PostgreSQLTarget) MigrationsExecuted() ([]time.Time, error) {
	migrations := make([]time.Time, 0, 10)
	err := target.withConn(func(ctx context.Context, conn *sql.Conn) error {
		rows, err := conn.QueryContext(ctx, fmt.Sprintf(`SELECT id FROM %s ORDER BY id`, target.tableName))
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var id time.Time
			if err := rows.Scan(&id); err != nil {
				return err
			}
			migrations = append(migrations, id)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return migrations, nil

}

func (target *PostgreSQLTarget) withConn(h func(ctx context.Context, conn *sql.Conn) error) error {
	ctx := context.Background()

	conn, err := target.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(context.Background(), fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id timestamptz NOT NULL PRIMARY KEY)", target.tableName))
	if err != nil {
		return err
	}

	return h(ctx, conn)
}
