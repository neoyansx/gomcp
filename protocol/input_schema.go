package protocol

import "sync"

type (
	Property struct {
		Name        string
		Type        string `json:"type"`
		Description string `json:"description,omitempty"`
	}
	InputSchema struct {
		sync.RWMutex
		Type       string              `json:"type"`
		Properties map[string]Property `json:"properties"`
	}
)

func newInputSchema() *InputSchema {
	return &InputSchema{
		Type:       "object",
		Properties: make(map[string]Property),
	}
}
