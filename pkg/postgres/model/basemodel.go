package model

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Make sure postgres is included in any build

	"github.com/fzerorubigd/engine/pkg/log"
)

var (
	dbmap *sqlx.DB
	db    *sql.DB
)

// Manager is a base manager for transaction model
type Manager struct {
	tx *sqlx.Tx

	transaction bool
}

// InTransaction return true if this manager s in transaction
func (m *Manager) InTransaction() bool {
	return m.transaction
}

// Begin is for begin transaction
func (m *Manager) Begin() error {
	var err error
	if m.transaction {
		log.Panic("Already in transaction")
	}
	m.tx, err = dbmap.Beginx()
	if err == nil {
		m.transaction = true
	}
	return err
}

// Commit is for committing transaction. panic if transaction is not started
func (m *Manager) Commit() error {
	if !m.transaction {
		log.Panic("Not in transaction")
	}
	err := m.tx.Commit()
	if err != nil {
		return err
	}
	m.tx = nil
	m.transaction = false
	return nil
}

// Rollback is for RollBack transaction. panic if transaction is not started
func (m *Manager) Rollback() error {
	if !m.transaction {
		log.Panic("Not in transaction")
	}
	err := m.tx.Rollback()

	if err != nil {
		return err
	}

	m.transaction = false
	return nil
}

// GetDbMap is for getting the current dbmap
func (m *Manager) GetDbMap() DBX {
	if m.transaction {
		return m.tx
	}
	return dbmap
}

// GetSQLDB return the raw connection to database
func (m *Manager) GetSQLDB() *sql.DB {
	return db
}

// GetDialect return the dialect of this instance
func (m *Manager) GetDialect() string {
	return "postgres"
}

// Hijack try to hijack into a transaction
func (m *Manager) Hijack(ts DBX) error {
	if m.transaction {
		return errors.New("already in transaction")
	}
	t, ok := ts.(*sqlx.Tx)
	if !ok {
		return errors.New("there is no transaction to hijack")
	}

	m.transaction = true
	m.tx = t

	return nil
}

// TruncateTables try to truncate tables , useful for tests
func (m *Manager) TruncateTables(cascade, resetIdentity bool, tbl ...string) error {
	q := "TRUNCATE " + strings.Join(tbl, " , ")
	if resetIdentity {
		q += " RESTART IDENTITY "
	}
	if cascade {
		q += " CASCADE "
	}

	_, err := dbmap.Exec(q)
	return err
}

// PrefixArray is helper function for appending table name to fields name
func (m *Manager) PrefixArray(prefix string, fields ...string) []string {
	ret := make([]string, len(fields))
	for i := range fields {
		ret[i] = prefix + fields[i]
	}

	return ret
}

// Initialize the module
func Initialize(d *sql.DB) {
	dbmap = sqlx.NewDb(d, "postgres")
	db = d
}
