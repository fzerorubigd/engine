package mockery

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"go.uber.org/zap/zaptest"

	"github.com/fzerorubigd/engine/pkg/config"
	"github.com/fzerorubigd/engine/pkg/initializer"
	"github.com/fzerorubigd/engine/pkg/log"
	"github.com/fzerorubigd/engine/pkg/postgres"
)

var (
	alreadyRegistered bool
)

// Start the mockery, used for tests only
func Start(ctx context.Context, t *testing.T) func() {
	if !alreadyRegistered {
		alreadyRegistered = true
		config.Initialize(ctx, "testing", "T")
		log.SwapLogger(zaptest.NewLogger(t))

		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			host.String(),
			port.Int(),
			user.String(),
			pass.String(),
			dbname.String(),
			sslmode.String(),
		)
		txdb.Register("txdb", "postgres", dsn)
		postgres.DefaultInitDB = sqltxTesting
	}

	return initializer.Initialize(ctx)
}
