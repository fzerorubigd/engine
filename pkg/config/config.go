package config

import (
	"runtime"
	"sync"

	"github.com/goraz/onion"
	_ "github.com/goraz/onion/loaders/yaml" // for loading yaml file

)

var (
	lock sync.RWMutex
	// map of conf,description
	descriptions = make(map[string]string)

	o = onion.New()
)

func initLayer() onion.Layer {
	data := map[string]interface{}{
		"core": map[string]interface{}{
			"env":               "dev",
			"max_cpu_available": runtime.NumCPU(),
			"machine_name":      "m1",
		},
	}
	setDescription("core.env", "environ to use")
	setDescription("core.max_cpu_available", "number of cpu")
	setDescription("core.machine_name", "machine name")

	return onion.NewMapLayer(data)
}
