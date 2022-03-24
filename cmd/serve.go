package cmd

import (
	"os"

	"github.com/hashicorp/go-hclog"
	hPlugin "github.com/hashicorp/go-plugin"
	"github.com/odpf/optimus/cmd/server"
	"github.com/odpf/optimus/config"
	"github.com/odpf/optimus/plugin"
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
	conf, err := config.LoadOptimusConfig()
	if err != nil {
		panic(err.Error())
	}

	// initiate jsonLogger
	jsonLogger := initLogger(jsonLoggerType, conf.Log)

	c.RunE = func(c *cli.Command, args []string) error {
		// initiate plugin log level
		pluginLogLevel := hclog.Info
		if conf.Log.Level == config.LogLevelDebug {
			pluginLogLevel = hclog.Debug
		}

		// discover and load plugins. TODO: refactor this
		if err := plugin.Initialize(hclog.New(&hclog.LoggerOptions{
			Name:   "optimus",
			Output: os.Stdout,
			Level:  pluginLogLevel,
		})); err != nil {
			return err
		}
		// Make sure we clean up any managed plugins at the end of this
		defer hPlugin.CleanupClients()

		// init telemetry
		teleShutdown, err := config.InitTelemetry(jsonLogger, conf.Telemetry)
		if err != nil {
			return err
		}
		defer teleShutdown()

		return server.Initialize(jsonLogger, *conf)
	}

	return c
}
