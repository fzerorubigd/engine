package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fzerorubigd/engine/cmd/qollenge"
	"github.com/fzerorubigd/engine/pkg/cli"
	"github.com/fzerorubigd/engine/pkg/initializer"
	"github.com/fzerorubigd/engine/pkg/log"
	"github.com/fzerorubigd/engine/pkg/migration"
	"github.com/fzerorubigd/engine/pkg/postgres/model"
)

var (
	action = flag.String("action", "up", "up/down is supported, default is up")
	n      int
)

func main() {
	ctx := cli.Context()
	qollenge.InitializeConfig(ctx)

	defer initializer.Initialize(ctx)()

	flag.Parse()
	var err error
	m := &model.Manager{}
	if *action == "up" {
		n, err = migration.Do(m, migration.Up, 0)
		if err != nil {
			log.Fatal("Migration failed", log.Err(err))
		}
		fmt.Printf("\n\n%d migration is applied\n", n)
	} else if *action == "down" {
		n, err = migration.Do(m, migration.Down, 1)
		if err != nil {
			log.Fatal("Migration failed", log.Err(err))
		}
		fmt.Printf("\n\n%d migration is applied\n", n)
	} else if *action == "down-all" {
		n, err = migration.Do(m, migration.Down, 0)
		if err != nil {
			log.Fatal("Migration failed", log.Err(err))
		}
		fmt.Printf("\n\n%d migration is applied\n", n)
	} else if *action == "redo" {
		n, err = migration.Do(m, migration.Down, 1)
		if err == nil {
			n, err = migration.Do(m, migration.Up, 1)
		}
		if err != nil {
			log.Fatal("Migration failed", log.Err(err))
		}
		fmt.Printf("\n\n%d migration is applied\n", n)

	} else if *action == "list" {
		migration.List(m, os.Stdout)
	}

	if err != nil {
		log.Fatal("Error on migration", log.Err(err))
	}
}
