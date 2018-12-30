package postgres

// TODO : multi connection support
import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/fzerorubigd/balloon/pkg/assert"
	"github.com/fzerorubigd/balloon/pkg/initializer"
	"github.com/fzerorubigd/balloon/pkg/log"
	"github.com/fzerorubigd/balloon/pkg/postgres/model"
	_ "github.com/lib/pq" // Make sure the pg is available
)

var (
	db  *sql.DB
	all []initializer.Simple

	// DefaultInitDB is a function for create a db instance, usable for change test db
	DefaultInitDB func(context.Context) (*sql.DB, error)
)

func realInstance(ctx context.Context) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host.String(),
		port.Int(),
		user.String(),
		pass.String(),
		dbname.String(),
		sslmode.String(),
	)
	pgDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	pgDB.SetMaxIdleConns(maxIdle.Int())
	pgDB.SetMaxOpenConns(maxCon.Int())

	// TODO : use better retry function
	cnt := 1
	for {
		err := pgDB.Ping()
		if err == nil {
			break
		}
		log.Error("Can not ping database", log.Err(err))
		select {
		case <-time.After(time.Second * time.Duration(cnt)):
			cnt++
		case <-ctx.Done():
			return nil, errors.New("context canceled")
		}
		if cnt > 10 {
			cnt = 10
		}
	}
	return pgDB, nil
}

// Hooker interface :))))) You have a dirty mind.
type Hooker interface {
	// AddHook is called after initialize only if the manager implement it
	AddHook()
}

type modelsInitializer struct {
}

func (modelsInitializer) Healthy(context.Context) error {
	return db.Ping()
}

// Initialize the modules, its safe to call this as many time as you want.
func (modelsInitializer) Initialize(ctx context.Context) {
	var err error
	db, err = DefaultInitDB(ctx)
	assert.Nil(err)
	model.Initialize(db)

	for i := range all {
		all[i].Initialize()

	}
	// If they are hooker call them.
	for i := range all {
		if h, ok := all[i].(Hooker); ok {
			h.AddHook()
		}
	}
	go func() {
		c := ctx.Done()
		if c == nil {
			return
		}
		<-c
		log.Debug("Postgres finalized")
	}()
	log.Debug("Postgres is ready")
}

// Register a new initializer module
func Register(m ...initializer.Simple) {
	all = append(all, m...)
}

func init() {
	DefaultInitDB = realInstance
	initializer.Register(&modelsInitializer{}, 0)
}
