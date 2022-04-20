package handlers

import (
	"cobra_starter/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/hibiken/asynq"
)

//---------------------------------------------------------------
// Write a function HandleXXXTask to handle the input task.
// Note that it satisfies the asynq.HandlerFunc interface.
//
// Handler doesn't need to be a function. You can define a type
// that satisfies asynq.Handler interface. See examples below.
//---------------------------------------------------------------

var count int

func FailHandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {

	var p models.EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	count++
	if math.Mod(float64(count), 1) == 0 {
		fmt.Printf("fail test run. enqueued task %d: Sending Email to User: user_id=%d, template_id=%s\n", count, p.UserID, p.TemplateID)
	}

	// Email delivery code ...
	return errors.New("failed to send email")
}
func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {

	var p models.EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	count++
	if math.Mod(float64(count), 1) == 0 {
		fmt.Printf("enqueued task %d: Sending Email to User: user_id=%d, template_id=%s\n", count, p.UserID, p.TemplateID)
	}

	// Email delivery code ...
	return nil
}
