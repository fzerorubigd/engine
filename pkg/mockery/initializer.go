package mockery

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/fzerorubigd/balloon/pkg/config"
	"github.com/fzerorubigd/balloon/pkg/initializer"
	"github.com/fzerorubigd/balloon/pkg/log"
	"github.com/fzerorubigd/balloon/pkg/postgres"
	"go.uber.org/zap/zaptest"
)

var (
	alreadyRegistered bool
)

func Start(t *testing.T, ctx context.Context) func() {
	if !alreadyRegistered {
		alreadyRegistered = true
		config.Initialize("testing", "T")
		log.SwappLogger(zaptest.NewLogger(t))

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
