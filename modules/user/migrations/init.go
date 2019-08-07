package migrations

import "github.com/fzerorubigd/engine/pkg/migration"

func init() {
	migration.Register(Asset, AssetDir, "postgres")
}
