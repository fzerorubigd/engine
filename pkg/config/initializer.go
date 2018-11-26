package config

import (
	"os"

	"github.com/fzerorubigd/balloon/pkg/log"

	"github.com/fzerorubigd/balloon/pkg/assert"
	"github.com/fzerorubigd/expand"
	onion "gopkg.in/fzerorubigd/onion.v3"
	"gopkg.in/fzerorubigd/onion.v3/extraenv"
)

var (
	all []Initializer
)

// Initializer is the config initializer for module
type Initializer interface {
	// Initialize is called when the module is going to add its layer
	Initialize() DescriptiveLayer
	// Loaded inform the modules that all layer are ready
	Loaded()
}

//Initialize try to initialize config
func Initialize(appName, prefix string, layers ...onion.Layer) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	assert.Nil(o.AddLayer(defaultLayer()))

	for i := range all {
		nL := all[i].Initialize()
		if nL != nil {
			_ = o.AddLayer(nL)
		}
	}

	// now add the layer provided by app
	for i := range layers {
		_ = o.AddLayer(layers[i])
	}
	// Now load external config to overwrite them all.
	if err := o.AddLayer(onion.NewFileLayer("/etc/" + appName + "/" + env + ".yaml")); err == nil {
		log.Info("Loading config", log.String("file", "/etc/"+appName+"/"+env+".yaml"))
	}
	p, err := expand.Path("$HOME/." + appName + "/" + env + ".yaml")
	if err == nil {
		if err = o.AddLayer(onion.NewFileLayer(p)); err == nil {
			log.Info("Loading config", log.String("file", p))
		}
	}

	p, err = expand.Path("$PWD/configs/" + appName + "/" + env + ".yaml")
	if err == nil {
		if err = o.AddLayer(onion.NewFileLayer(p)); err == nil {
			log.Info("Loading config", log.String("file", p))
		}
	}

	o.AddLazyLayer(extraenv.NewExtraEnvLayer(prefix))

	// load all registered variables
	o.Load()
	// tell them that every thing is loaded
	for i := range all {
		all[i].Loaded()
	}
}

func setDescription(key, desc string) {
	lock.Lock()
	defer lock.Unlock()
	if d, ok := configs[key]; ok && d != "" && desc == "" {
		// if the new description is empty and the old one is not, ignore the new one
		return
	}
	configs[key] = desc
}

// Register a config module
func Register(i ...Initializer) {
	all = append(all, i...)
}
