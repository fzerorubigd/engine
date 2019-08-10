package migrations

import "elbix.dev/engine/pkg/migration"

func init() {
	migration.Register(Asset, AssetDir, "postgres")
}
