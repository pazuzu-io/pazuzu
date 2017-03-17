package project

import (
	"github.com/urfave/cli"
	"github.com/zalando-incubator/pazuzu/cli/pazuzu/actions"
)

var Command = cli.Command{
	Name:  "project",
	Usage: "Configure Pazuzu project settings",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "d, directory",
			Usage: "Sets source path where project configuration is located.",
		},
	},
	Subcommands: []cli.Command{
		{
			Name:   "add",
			Usage:  "Add feature to the project",
			Action: actions.ProjectAddFeatures,
		},
		{
			Name:   "list",
			Usage:  "List features used in the project",
			Action: actions.ProjectListFeatures,
		},
		{
			Name:   "remove",
			Usage:  "Remove feature from the project",
			Action: actions.ProjectRemoveFeatures,
		},
		{
			Name:   "clean",
			Usage:  "Remove all Pazuzu-generated files",
			Action: actions.ProjectClean,
		},
		{
			Name:  "build",
			Usage: "Build and test Docker image based on the project configuration",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "n, name",
					Usage: "Set the name for Docker image",
				},
			},
			Action: actions.ProjectBuild,
		},
		{
			Name:   "show",
			Usage:  "Show base image, author, license settings of the project",
			Action: actions.ProjectShow,
		},
		{
			Name:   "set",
			Usage:  "Set base image, author, license settings of the project",
			Action: actions.ProjectSet,
		},
	},
}
