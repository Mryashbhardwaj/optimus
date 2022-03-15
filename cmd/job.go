package cmd

import (
	"github.com/odpf/optimus/config"

	"github.com/odpf/optimus/models"
	"github.com/odpf/salt/log"
	cli "github.com/spf13/cobra"
)

func jobCommand(l log.Logger, conf config.Optimus, pluginRepo models.PluginRepository) *cli.Command {
	cmd := &cli.Command{
		Use:   "job",
		Short: "Interact with schedulable Job",
		Annotations: map[string]string{
			"group:core": "true",
		},
	}

	cmd.AddCommand(jobCreateCommand(l, conf, pluginRepo))
	cmd.AddCommand(jobAddHookCommand(l, conf, pluginRepo))
	cmd.AddCommand(jobRenderTemplateCommand(l, conf, pluginRepo))
	cmd.AddCommand(jobValidateCommand(l, conf, pluginRepo, conf.Project.Name, conf.Host))
	cmd.AddCommand(jobRunCommand(l, conf, pluginRepo, conf.Project.Name, conf.Host))
	cmd.AddCommand(jobStatusCommand(l, conf.Project.Name, conf.Host))
	cmd.AddCommand(jobRefreshCommand(l, conf))
	return cmd
}
