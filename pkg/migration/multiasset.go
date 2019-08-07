package migration

import (
	"database/sql"
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/rubenv/sql-migrate"

	// Make sure postgres is already initialized
	_ "github.com/fzerorubigd/engine/pkg/postgres"
)

var (
	all multiAsset
)

// Manager is the simple model manager
type Manager interface {
	// Get the actual sql.db pool
	GetSQLDB() *sql.DB
	// GetDialect must return the dialect
	GetDialect() string
}

// Direction is a helper to prevent import the migration library in other packages
type Direction int

const (
	// Up means migrate up
	Up Direction = iota
	// Down means migrate down
	Down
)

type byID []*migrate.Migration

func (b byID) Len() int           { return len(b) }
func (b byID) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b byID) Less(i, j int) bool { return b[i].Less(b[j]) }

type multiAsset []migrate.AssetMigrationSource

func (m multiAsset) FindMigrations() ([]*migrate.Migration, error) {
	var mig []*migrate.Migration
	for i := range m {
		m, err := m[i].FindMigrations()
		if err != nil {
			return nil, err
		}
		mig = append(mig, m...)
	}
	sort.Sort(byID(mig))
	return mig, nil
}

// Register a new asset function
func Register(asset func(path string) ([]byte, error), assetDir func(path string) ([]string, error), dir string) {
	all = append(all, migrate.AssetMigrationSource{
		Asset:    asset,
		AssetDir: assetDir,
		Dir:      dir,
	})
}

// Do is my try to migrate on demand. but I don't know if there is more than
// one instance is in memory and if that make trouble
func Do(db Manager, dir Direction, max int) (int, error) {
	var (
		err error
		n   int
	)
	if len(all) == 0 {
		return 0, nil
	}
	if max == 0 {
		n, err = migrate.Exec(db.GetSQLDB(), db.GetDialect(), all, migrate.MigrationDirection(dir))
	} else {
		n, err = migrate.ExecMax(db.GetSQLDB(), db.GetDialect(), all, migrate.MigrationDirection(dir), max)
	}
	if err != nil {
		return 0, err
	}

	return n, nil
}

// List already applied migration and format them in io writer
func List(db Manager, target io.Writer) {
	mig, _ := migrate.GetMigrationRecords(db.GetSQLDB(), db.GetDialect())
	w := tabwriter.NewWriter(target, 0, 8, 0, '\t', 0)
	_, _ = fmt.Fprintln(w, "|ID\t|Applied at\t|")
	for i := range mig {
		_, _ = fmt.Fprintf(w, "|%s\t|%s\t|\n", mig[i].Id, mig[i].AppliedAt)
	}
	_ = w.Flush()
}
