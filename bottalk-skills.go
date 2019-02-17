package main

import (
	"github.com/urfave/cli"
	"log"
)

func getSkillCommands() []cli.Command {

	return []cli.Command{

		{

			Name:  "list",
			Usage: "Display list of your skills",
			Action: func(c *cli.Context) error {
				log.Println("Fetching list of your skills")
				getSkillList()
				return nil
			},
		},
		{

			Name:  "new",
			Usage: "Create new skill",
			Action: func(c *cli.Context) error {
				skillLanguage := c.Args().Get(1)
				skillName := c.Args().Get(0)
				createNewSkill(skillName, skillLanguage)
				return nil
			},
		},
		{
			Name:  "pull",
			Usage: "Pull skill and write files into new folder",
			Action: func(c *cli.Context) error {
				skillToken := c.Args().Get(0)
				log.Println("Requesting skill " + skillToken)
				getSkillFiles(skillToken)
				return nil
			},
		}, {
			Name:  "push",
			Usage: "Push skill files from current folder to bottalk server",
			Action: func(c *cli.Context) error {
				pushSkillFiles()
				return nil
			},
		}}
}
