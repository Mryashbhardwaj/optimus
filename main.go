package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	_ "net/http/pprof"

	"github.com/odpf/optimus/cmd"
	_ "github.com/odpf/optimus/ext/datastore"
	"github.com/odpf/optimus/models"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	command := cmd.New(
		models.PluginRegistry,
		models.DatastoreRegistry,
	)

	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
