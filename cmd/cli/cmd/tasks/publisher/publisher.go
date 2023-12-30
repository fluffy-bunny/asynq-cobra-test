/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package publisher

import (
	"cobra_starter/internal/cobra_utils"
	"cobra_starter/internal/models"
	"fmt"
	"math"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var count int
var queuename string
var fail bool
var group string

// publisherCmd represents the publisher command
var publisherCmd = &cobra.Command{
	Use:               "publisher",
	PersistentPreRunE: cobra_utils.ParentPersistentPreRunE,
	Short:             "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log := zerolog.Ctx(cmd.Context()).With().Logger()
		client := asynq.NewClient(asynq.RedisClientOpt{Addr: "localhost:6379",
			Network:  "tcp",
			Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81"})
		defer client.Close()
		var options []asynq.Option
		if queuename != "" {
			options = append(options, asynq.Queue(queuename))
		}
		if len(group) > 0 {
			options = append(options, asynq.Group(group))
		}
		log.Info().Interface("options", options).Msg("options")
		for i := 0; i < count; i++ {
			task, err := models.NewEmailDeliveryTask(i, "some:template:id", fail)
			if err != nil {
				log.Fatal().Err(err).Msg("could not create task")
			}
			info, err := client.Enqueue(task, options...)
			if err != nil {
				log.Fatal().Err(err).Msg("could not enqueue task")
			}
			if math.Mod(float64(i), 1000) == 0 {
				fmt.Printf("enqueued task %d: %s\n", i, info.ID)
			}

		}
		fmt.Printf("enqueued %d tasks fail:%v\n", count, fail)

	},
}

func InitCommand(parent *cobra.Command) {
	parent.AddCommand(publisherCmd)
	publisherCmd.Flags().IntVarP(&count, "count", "c", 1, "count")
	publisherCmd.Flags().BoolVarP(&fail, "fail", "f", false, "fail to handle the task")
	publisherCmd.Flags().StringVarP(&group, "group", "g", "", "enque to an aggregate group")

	publisherCmd.Flags().StringVarP(&queuename, "queuename", "q", "", "queuename")
}
