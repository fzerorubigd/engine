package migrations

import "github.com/fzerorubigd/balloon/pkg/migration"

func init() {
	migration.Register(Asset, AssetDir, "postgres")
}
