package events

import (
	"fmt"

	"github.com/google/uuid"
)

type Events string

type Metadata struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type BaseEvent struct {
	Metadata Metadata `json:"metadata"`
}

func (self *BaseEvent) GenerateEventId(event_type string, user_id string) string {
	random_id := uuid.New().String()
	idempotency := fmt.Sprintf("%v-%v-%v", event_type, user_id, random_id)

	self.Metadata.Id = idempotency
	return idempotency
}
