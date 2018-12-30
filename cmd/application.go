package cmd

import (
	// enable accounting module
	_ "github.com/fzerorubigd/balloon/modules/accounting"
	// enable misc module
	_ "github.com/fzerorubigd/balloon/modules/misc"
	// enable the user module
	_ "github.com/fzerorubigd/balloon/modules/user"
	"github.com/fzerorubigd/balloon/pkg/config"
)

// InitializeConfig for this application
func InitializeConfig() {
	config.Initialize("balloon", "BAL")
}
