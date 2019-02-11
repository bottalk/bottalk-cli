package main

import (
	"github.com/urfave/cli"
	"log"
)

func getSkillCommands() []cli.Command {

	return []cli.Command{{

		Name:  "skill",
		Usage: "Skill related commands",
		Subcommands: []cli.Command{
			{

				Name:  "list",
				Usage: "Display list of your skills",
				Action: func(c *cli.Context) error {
					log.Println("Fetching list of your skills")
					getSkillList()
					return nil
				},
			},
		},
	}, {
		Name:  "pull",
		Usage: "pull skill and write files into disk",
		Action: func(c *cli.Context) error {
			skillToken := c.Args().Get(0)
			log.Println("Requesting skill " + skillToken)
			getSkillFiles(skillToken)
			return nil
		},
	}}
}
