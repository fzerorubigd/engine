package cmd

import (
	_ "github.com/fzerorubigd/balloon/modules/misc/impl"
	_ "github.com/fzerorubigd/balloon/modules/user/impl"
	"github.com/fzerorubigd/balloon/pkg/config"
)

// InitializeConfig for this application
func InitializeConfig() {
	config.Initialize("balloon", "BAL")
}
