package events

import (
	"fmt"

	"github.com/google/uuid"
)

type Events string
type Publisher func(data any) (any, error)

type Metadata struct {
	Id    string `json:"id"`
	Reply string `json:"reply"`
}

// Generic type for each event
type BaseEvent struct {
	Metadata Metadata `json:"metadata"`
}

func (self *BaseEvent) GenerateEventId(event_type string, user_id string) string {
	random_id := uuid.New().String()
	idempotency := fmt.Sprintf("%v-%v-%v", event_type, user_id, random_id)

	self.Metadata.Id = idempotency
	return idempotency
}
