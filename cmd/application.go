package cmd

import (
	_ "github.com/fzerorubigd/balloon/modules/accounting"
	_ "github.com/fzerorubigd/balloon/modules/misc"
	_ "github.com/fzerorubigd/balloon/modules/user"
	"github.com/fzerorubigd/balloon/pkg/config"
)

// InitializeConfig for this application
func InitializeConfig() {
	config.Initialize("balloon", "BAL")
}
