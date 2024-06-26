package database

import (
	"context"
	"database/sql"
	"io"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB                    *sqlx.DB
	url                   string
	maxOpenConnections    int
	maxIdleConnections    int
	connectionMaxLifetime time.Duration
	connectionMaxIdleTime time.Duration
	log                   *log.Logger
}

type NewDatabaseOptions struct {
	URL                   string
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
	ConnectionMaxIdleTime time.Duration
	Log                   *log.Logger
}

// NewDatabase with the given options.
// If no logger is provided, logs are discarded.
func NewDatabase(opts NewDatabaseOptions) *Database {
	if opts.Log == nil {
		opts.Log = log.New(io.Discard, "", 0)
	}

	// - Set WAL mode (not strictly necessary each time because it's persisted in the database, but good for first run)
	// - Set busy timeout, so concurrent writers wait on each other instead of erroring immediately
	// - Enable foreign key checks
	opts.URL += "?_journal=WAL&_timeout=5000&_fk=true"

	return &Database{
		url:                   opts.URL,
		maxOpenConnections:    opts.MaxOpenConnections,
		maxIdleConnections:    opts.MaxIdleConnections,
		connectionMaxLifetime: opts.ConnectionMaxLifetime,
		connectionMaxIdleTime: opts.ConnectionMaxIdleTime,
		log:                   opts.Log,
	}
}

func (d *Database) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d.log.Println("Connecting to database at", d.url)

	var err error
	d.DB, err = sqlx.ConnectContext(ctx, "sqlite3", d.url)
	if err != nil {
		return err
	}

	d.log.Println("Setting connection pool options (",
		"max open connections:", d.maxOpenConnections,
		", max idle connections:", d.maxIdleConnections,
		", connection max lifetime:", d.connectionMaxLifetime,
		", connection max idle time:", d.connectionMaxIdleTime,
		")")
	d.DB.SetMaxOpenConns(d.maxOpenConnections)
	d.DB.SetMaxIdleConns(d.maxIdleConnections)
	d.DB.SetConnMaxLifetime(d.connectionMaxLifetime)
	d.DB.SetConnMaxIdleTime(d.connectionMaxIdleTime)

	return nil
}

type Transaction struct {
	Tx *sqlx.Tx
}

func (d *Database) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Transaction, error) {
	tx, err := d.DB.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &Transaction{Tx: tx}, nil
}

func (t *Transaction) Commit() error {
	return t.Tx.Commit()
}

func (t *Transaction) Rollback() error {
	return t.Tx.Rollback()
}

// func CheckForLocks(db *sqlx.DB) error {
//     type LockStatus struct {
//         Database string `db:"database"`
//         Table    string `db:"table"`
//         Type     string `db:"type"`
//     }
//
//     var locks []LockStatus
//     err := db.Select(&locks, "PRAGMA lock_status")
//     if err != nil {
//         return err
//     }
//
//     if len(locks) == 0 {
//         log.Println("No locks found.")
//     } else {
//         for _, lock := range locks {
//             log.Printf("Lock found: Database=%s, Table=%s, Type=%s\n", lock.Database, lock.Table, lock.Type)
//         }
//     }
//
//     return nil
// }
