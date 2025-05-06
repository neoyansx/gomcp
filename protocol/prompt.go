package protocol

type (
	Prompt struct {
		Name        string            `json:"name"`
		Description string            `json:"description,omitempty"`
		Args        []*PromptArgument `json:"arguments,omitempty"`
		handler     HandleFunc
		role        string
		mimeType    string
		contentType ContentType
	}

	PromptArgument struct {
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
		IsRequired  bool   `json:"required,omitempty"`
	}
)

func (p *Prompt) getMime() string {
	return p.mimeType
}

func (p *Prompt) invoke(arguments map[string]any) (string, error) {
	return p.handler(arguments)
}

func (p *Prompt) getRole() string {
	return p.role
}

func (p *Prompt) getDescription() string {
	return p.Description
}

func (p *Prompt) getContent() ContentType {
	return p.contentType
}

// Content mark content type of handler result
func (p *Prompt) Content(typ ContentType) IPrompt {
	if typ == Image || typ == Audio || typ == Text {
		p.contentType = typ
	} else {
		p.contentType = Text
	}
	return p
}

func (p *Prompt) Describe(description string) IPrompt {
	p.Description = description
	return p
}

func (p *Prompt) Handler(handler HandleFunc) IPrompt {
	p.handler = handler
	return p
}

func (p *Prompt) SetDescribe(description string) {
	p.Description = description
}

func (p *Prompt) SetHandler(handler HandleFunc) {
	p.handler = handler
}

func (p *Prompt) Role(role string) IPrompt {
	p.role = role
	return p
}

// Argument add argument without description to Prompt directly
func (p *Prompt) Argument(name string, required bool) IPrompt {
	arg := &PromptArgument{
		Name:       name,
		IsRequired: required,
	}
	p.Args = append(p.Args, arg)
	return p
}

func (p *Prompt) getName() string {
	return p.Name
}
