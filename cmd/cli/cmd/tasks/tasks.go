/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package tasks

import (
	"fmt"

	"cobra_starter/cmd/cli/shared"
	"cobra_starter/internal/cobra_utils"

	"cobra_starter/cmd/cli/cmd/tasks/handler"
	"cobra_starter/cmd/cli/cmd/tasks/publisher"

	"github.com/spf13/cobra"
)

// tasksCmd represents the tasks command
var tasksCmd = &cobra.Command{
	Use:               "tasks",
	PersistentPreRunE: cobra_utils.ParentPersistentPreRunE,
	Short:             "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %s\n", shared.Version)
	},
}

func InitCommand(parent *cobra.Command) {
	parent.AddCommand(tasksCmd)
	handler.InitCommand(tasksCmd)
	publisher.InitCommand(tasksCmd)
}
