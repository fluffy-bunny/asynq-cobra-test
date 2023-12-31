/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package handler

import (
	"cobra_starter/internal/cobra_utils"
	task_handlers "cobra_starter/internal/handlers"
	"cobra_starter/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"

	asynq_log "cobra_starter/internal/log"

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
		log := zerolog.Ctx(cmd.Context()).With().Logger()
		queuesMap := make(map[string]int)
		for _, queue := range queues {
			items := strings.Split(queue, ":")
			if len(items) != 2 {
				log.Fatal().Msgf("invalid queue: %s", queue)
			}
			key := items[0]
			value, err := strconv.Atoi(items[1])
			if err != nil {
				log.Fatal().Msgf("invalid queue: %s", queue)
			}
			queuesMap[key] = value

		}
		aLog := asynq_log.NewLogger()
		cfg := asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: queuesMap,
			// See the godoc for other configuration options
			GroupAggregator:  asynq.GroupAggregatorFunc(aggregateHandler),
			GroupMaxDelay:    1 * time.Minute,
			GroupGracePeriod: 30 * time.Second,
			GroupMaxSize:     500,
			BaseContext: func() context.Context {
				return cmd.Context()
			},
			Logger: aLog,
			IsFailure: func(err error) bool {

				return true
			},
			ErrorHandler: asynq.ErrorHandlerFunc(reportError),
		}

		srv := asynq.NewServer(
			asynq.RedisClientOpt{Addr: "localhost:6379",
				Network:  "tcp",
				Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81"},
			cfg,
		)

		// mux maps a type to a handler
		mux := asynq.NewServeMux()
		mux.HandleFunc(models.TypeEmailDelivery, task_handlers.HandleEmailDeliveryTask)
		mux.HandleFunc("aggregated-task", handleAggregatedTask)

		// ...register other handlers...
		sb := strings.Builder{}
		for k, v := range queuesMap {
			sb.WriteString(fmt.Sprintf("%s:%d, ", k, v))
		}
		fmt.Printf("Listening to queues: %s\n", sb.String())
		if err := srv.Run(mux); err != nil {
			log.Fatal().Err(err).Msg("could not run server")
		}
	},
}

func reportError(ctx context.Context, task *asynq.Task, err error) {
	log := zerolog.Ctx(ctx).With().Logger()
	retried, _ := asynq.GetRetryCount(ctx)
	maxRetry, _ := asynq.GetMaxRetry(ctx)
	log = log.With().Int("retried", retried).Int("maxRetry", maxRetry).Logger()

	if retried >= maxRetry {
		log.Error().Err(err).Msg("error handling task")
	}
}

func InitCommand(parent *cobra.Command) {
	parent.AddCommand(handlerCmd)

	handlerCmd.Flags().StringArrayVarP(&queues, "queues", "q", []string{"critical:6", "default:3", "low:1", "aggregate:3"}, "queues to listen to")
	handlerCmd.Flags().BoolVarP(&fail, "fail", "f", false, "fail to handle the task")

}

type aggregatedPayloads struct {
	Group    string
	Payloads [][]byte
}

// This function is used to aggregate multiple tasks into one.
func aggregateHandler(group string, tasks []*asynq.Task) *asynq.Task {
	// ... Your logic to aggregate the given tasks and return the aggregated task.
	// ... Use NewTask(typename, payload, opts...) to create a new task and set options if needed.
	// ... (Note) Queue option will be ignored and the aggregated task will always be enqueued to

	final := &aggregatedPayloads{
		Group: group,
	}
	for _, t := range tasks {
		final.Payloads = append(final.Payloads, t.Payload())
	}
	d, err := json.Marshal(final)
	if err != nil {
		panic(err)
	}
	var b strings.Builder
	for _, t := range tasks {
		b.Write(t.Payload())
		b.WriteString("\n")
	}
	return asynq.NewTask("aggregated-task", d)
}

func handleAggregatedTask(ctx context.Context, task *asynq.Task) error {
	var err error
	log := zerolog.Ctx(ctx).With().Logger()
	log.Info().Msg("Handler received aggregated task")
	final := &aggregatedPayloads{}
	groupS := &models.Group{}
	err = json.Unmarshal(task.Payload(), &final)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(final.Group), groupS)
	if err != nil {
		return err
	}
	// create a temporary struct to unmarshal the payload
	type summary struct {
		Group *models.Group
		Total int
	}
	summaryV := &summary{
		Group: groupS,
		Total: len(final.Payloads),
	}

	// you have to eat the error here so that this task doesn't get into a retry loop and finally archived.
	// t
	log.Info().Interface("summary", summaryV).Msg("unmarshalling payload")
	return nil
}
