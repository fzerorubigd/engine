package config

import (
	"runtime"
	"sync"

	"github.com/fzerorubigd/balloon/pkg/assert"
	"gopkg.in/fzerorubigd/onion.v3"
	_ "gopkg.in/fzerorubigd/onion.v3/yamlloader" // for loading yaml file
)

var (
	// map of conf,description
	configs = make(map[string]string)
	locker  = &sync.Mutex{}
	o       = onion.New()
)

// DescriptiveLayer is based on onion layer interface
type DescriptiveLayer interface {
	onion.Layer
	// Add get Description, key and value
	Add(string, string, interface{})
}

func defaultLayer() onion.Layer {
	d := layer{}
	d.Add("core.env", "dev", "environ to use")
	d.Add("core.max_cpu_available", runtime.NumCPU(), "number of cpu")
	d.Add("core.machine_name", "m1", "machine name")
	return &d
}

// layer is configuration holder
type layer struct {
	onion.DefaultLayer
}

// Load a layer into the Onion. the call is only done in the
// registration
func (l *layer) Load() (map[string]interface{}, error) {
	if l.DefaultLayer != nil {
		return l.DefaultLayer.Load()
	}

	return map[string]interface{}{}, nil
}

// Add set a default value for a key
func (l *layer) Add(key string, value interface{}, description string) {
	locker.Lock()
	defer locker.Unlock()
	if l.DefaultLayer == nil {
		l.DefaultLayer = onion.NewDefaultLayer()
	}
	assert.Nil(l.DefaultLayer.SetDefault(key, value))
	setDescription(key, description)

}

// GetDescriptions return config key, description
func GetDescriptions() map[string]string {
	return configs
}
