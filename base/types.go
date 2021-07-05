package base

import (
	"context"
	"encoding/json"
)

// EyewaEvent a base representation of an event fired/received
type EyewaEvent struct {
	ID        string `json:"id"`         // can be used for tracing
	Name      string `json:"name"`       // name of event - ProductUpdated, ProductDeleted etc
	EventType string `json:"event_type"` // type of event's entity - Product, Order etc

	// a representation on an error. provides reasons when a message ends up back
	// in the queue
	Errors []Error `json:"errors" binding:"omitempty"`

	// actual event payload
	Payload json.RawMessage `json:"payload"`

	// ts of when event was created
	CreatedAt string `json:"created_at"` // time in RFC3339 format
}

// Error a structural info about an error within the ecosystem
type Error struct {
	ErrorCode    int    `json:"error_code"`    // custom or http code should suffice
	ErrorMessage string `json:"error_message"` // error being reported
	CreatedAt    string `json:"created_at"`    // time in RFC3339 format
}

// EyewaEventError a structural error info about an event that failed
// during processing by a consumer - usecase: for publishing to a deadletter queue
// It is paramount to keep records of an event that failed and why.
type EyewaEventError struct {
	Event        string `json:"event"`         // string representation of event that was consumed off the queue but failed.
	ErrorMessage string `json:"error_message"` // error being reported
	CreatedAt    string `json:"created_at"`    // time in RFC3339 format
}

// MessageBrokerCallbackFunc all broker clients should define this callback fn
// so as to react to the state of events published/consumed - success/failure
type MessageBrokerCallbackFunc func(ctx context.Context, event *EyewaEvent, err error) error
