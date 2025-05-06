package types

import (
	"errors"
	"slices"
)

var levels = []string{"debug", "info", "notice", "warning", "error", "critical", "alert", "emergency"}

type Logging struct {
	Level  string `json:"level"`
	Logger string `json:"logger,omitempty"`
	Data   any    `json:"data,omitempty"`
}

func (l *Logging) SetLevel(level string) error {
	if !slices.Contains(levels, level) {
		return errors.New("invalid log level")
	}
	l.Level = level
	return nil
}
