/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package handler

import (
	"cobra_starter/internal/cobra_utils"
	task_handlers "cobra_starter/internal/handlers"
	"cobra_starter/internal/models"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/spf13/cobra"
)

var queues []string
var fail bool

// handlerCmd represents the handler command
var handlerCmd = &cobra.Command{
	Use:               "handler",
	PersistentPreRunE: cobra_utils.ParentPersistentPreRunE,
	Short:             "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		queuesMap := make(map[string]int)
		for _, queue := range queues {
			items := strings.Split(queue, ":")
			if len(items) != 2 {
				log.Fatalf("invalid queue: %s", queue)
			}
			key := items[0]
			value, err := strconv.Atoi(items[1])
			if err != nil {
				log.Fatalf("invalid queue: %s", queue)
			}
			queuesMap[key] = value

		}

		srv := asynq.NewServer(
			asynq.RedisClientOpt{Addr: "localhost:6379",
				Network:  "tcp",
				Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81"},
			asynq.Config{
				// Specify how many concurrent workers to use
				Concurrency: 10,
				// Optionally specify multiple queues with different priority.
				Queues: queuesMap,
				// See the godoc for other configuration options
			},
		)

		// mux maps a type to a handler
		mux := asynq.NewServeMux()
		if fail {
			mux.HandleFunc(models.TypeEmailDelivery, task_handlers.FailHandleEmailDeliveryTask)
		} else {
			mux.HandleFunc(models.TypeEmailDelivery, task_handlers.HandleEmailDeliveryTask)
		}
		// ...register other handlers...
		sb := strings.Builder{}
		for k, v := range queuesMap {
			sb.WriteString(fmt.Sprintf("%s:%d, ", k, v))
		}
		fmt.Printf("Listening to queues: %s\n", sb.String())
		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	},
}

func InitCommand(parent *cobra.Command) {
	parent.AddCommand(handlerCmd)

	handlerCmd.Flags().StringArrayVarP(&queues, "queues", "q", []string{"critical:6", "default:3", "low:1"}, "queues to listen to")
	handlerCmd.Flags().BoolVarP(&fail, "fail", "f", false, "fail to handle the task")
}
