package protocol

type PromptRole uint

const (
	User PromptRole = iota
	System
)

var PromptRoles = map[PromptRole]string{
	User:   "user",
	System: "system",
}

type (
	messages []*Message
	Message  struct {
		Ro      string   `json:"role"`
		Content IContent `json:"content"`
	}
)

func NewPromptMessage(role string) *Message {
	return &Message{
		Ro: role,
	}
}

func (m *Message) Role(role string) *Message {
	m.Ro = role
	return m
}

func (m *Message) AddContent(content IContent) {
	m.Content = content
}

func (ms *messages) addMessage(msg *Message) {
	*ms = append(*ms, msg)
}
