package config

import (
	"context"
	"os"

	"github.com/fzerorubigd/expand"
	"github.com/goraz/onion"
	"github.com/goraz/onion/configwatch"

	"github.com/fzerorubigd/engine/pkg/log"
)

// Initialize try to initialize config
func Initialize(ctx context.Context, appName, prefix string, layers ...onion.Layer) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	o.AddLayers(layers...)
	// Now load external config to overwrite them all.
	if l, err := onion.NewFileLayer("/etc/"+appName+"/"+env+".yaml", nil); err == nil {
		log.Info("Loading config", log.String("file", "/etc/"+appName+"/"+env+".yaml"))
		o.AddLayers(l)
	}
	p, err := expand.Path("$HOME/." + appName + "/" + env + ".yaml")
	if err == nil {
		if l, err := onion.NewFileLayer(p, nil); err == nil {
			log.Info("Loading config", log.String("file", p))
			o.AddLayers(l)
		}
	}

	p, err = expand.Path("$PWD/configs/" + appName + "/" + env + ".yaml")
	if err == nil {
		if l, err := onion.NewFileLayer(p, nil); err == nil {
			log.Info("Loading config", log.String("file", p))
			o.AddLayers(l)
		}
	}

	n := onion.NewEnvLayerPrefix("_", prefix)
	o.AddLayers(n)
	configwatch.WatchContext(ctx, o)
}

func setDescription(key, desc string) {
	lock.Lock()
	defer lock.Unlock()
	if d, ok := descriptions[key]; ok && d != "" && desc == "" {
		// if the new description is empty and the old one is not, ignore the new one
		return
	}
	descriptions[key] = desc
}
