package models

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeEmailDelivery = "email:deliver"
)

type EmailDeliveryPayload struct {
	UserID     int
	TemplateID string
	// marks this as a message that all the workers will fail to process
	FailIt bool
}

//----------------------------------------------
// Write a function NewXXXTask to create a task.
// A task consists of a type and a payload.
//----------------------------------------------

func NewEmailDeliveryTask(userID int, tmplID string, failit bool) (*asynq.Task, error) {
	payload, err := json.Marshal(EmailDeliveryPayload{
		UserID:     userID,
		TemplateID: tmplID,
		FailIt:     failit,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeEmailDelivery, payload), nil
}
