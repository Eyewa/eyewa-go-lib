package base

import (
	"encoding/json"
	"time"
)

// EyewaEvent a base representation of an event either fired/received
type EyewaEvent struct {
	ID        string `json:"id"`         // can be used for tracing
	Name      string `json:"name"`       // name of event - ProductUpdated, ProductDeleted etc
	EventType string `json:"event_type"` // type of event's entity - Product, Order etc

	// a representation on an error. provides reasons when a message ends up back
	// in the queue
	Error []Error `json:"error" binding:"omitempty"`

	// actual event payload
	Payload json.RawMessage `json:"payload"`

	// ts of when event was created
	CreatedAt string `json:"created_at"` // time in RFC3339 format
}

// Error a structural info about an error within the ecosystem
type Error struct {
	ErrorCode int        `json:"error_code"` // custom or http code should suffice
	Message   string     `json:"message"`    // error being reported
	CreatedAt *time.Time `json:"created_at"`
}
