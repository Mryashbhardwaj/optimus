package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
	hPlugin "github.com/hashicorp/go-plugin"
	"github.com/odpf/optimus/cmd/server"
	"github.com/odpf/optimus/config"
	"github.com/odpf/optimus/plugin"
	"github.com/odpf/salt/log"
	cli "github.com/spf13/cobra"
)

func serveCommand() *cli.Command {
	c := &cli.Command{
		Use:     "serve",
		Short:   "Starts optimus service",
		Example: "optimus serve",
		Annotations: map[string]string{
			"group:other": "dev",
		},
	}

	// TODO: find a way to load the config in one place
	conf, err := config.LoadServerConfig()
	if err != nil {
		panic(err.Error())
		return nil
	}

	// initiate jsonLogger
	var jsonLogger log.Logger
	pluginLogLevel := hclog.Info
	if conf.Log.Level != "" {
		jsonLogger = log.NewLogrus(log.LogrusWithLevel(conf.Log.Level), log.LogrusWithWriter(os.Stderr))
		if strings.ToLower(conf.Log.Level) == "debug" {
			pluginLogLevel = hclog.Debug
		}
	} else {
		jsonLogger = log.NewLogrus(log.LogrusWithLevel("INFO"), log.LogrusWithWriter(os.Stderr))
	}

	// discover and load plugins. TODO: refactor this
	if err := plugin.Initialize(hclog.New(&hclog.LoggerOptions{
		Name:   "optimus",
		Output: os.Stdout,
		Level:  pluginLogLevel,
	})); err != nil {
		hPlugin.CleanupClients()
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
	// Make sure we clean up any managed plugins at the end of this
	defer hPlugin.CleanupClients()

	// init telemetry
	teleShutdown, err := config.InitTelemetry(jsonLogger, conf.Telemetry)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
	defer teleShutdown()

	c.RunE = func(c *cli.Command, args []string) error {
		return server.Initialize(l, *conf)
	}

	return c
}
