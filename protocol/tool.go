package protocol

type (
	Tool struct {
		Name        string       `json:"name"`
		Description string       `json:"description"`
		InputSchema *InputSchema `json:"inputSchema,omitempty"`
		handler     HandleFunc
		contentType ContentType // handler返回结果的类型
		mimeType    string
		//resultType  string
	}

	CallToolRequest struct {
		//BasicMessage
		Params *Params `json:"params"`
	}

	Params struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}
)

func (t *Tool) getMime() string {
	return t.mimeType
}

func (t *Tool) getName() string {
	return t.Name
}

func (t *Tool) getContent() ContentType {
	return t.contentType
}

// Content mark type of handler result
func (t *Tool) Content(typ ContentType) ITool {
	if typ == Image || typ == Audio || typ == Text {
		t.contentType = typ
	} else {
		t.contentType = Text
	}
	return t
}

func (t *Tool) Describe(description string) ITool {
	t.Description = description
	return t
}

func (t *Tool) Handler(handler HandleFunc) ITool {
	t.handler = handler
	return t
}

func (t *Tool) SetDescribe(description string) {
	t.Description = description
}

func (t *Tool) SetHandler(handler HandleFunc) {
	t.handler = handler
}

func (t *Tool) Property(property Property) ITool {
	if t.InputSchema == nil {
		t.InputSchema = newInputSchema()
	}
	t.InputSchema.Lock()
	defer t.InputSchema.Unlock()
	t.InputSchema.Properties[property.Name] = property
	return t
}

func (t *Tool) Properties(ps ...Property) ITool {
	if t.InputSchema == nil {
		t.InputSchema = newInputSchema()
	}
	t.InputSchema.Lock()
	defer t.InputSchema.Unlock()
	for _, property := range ps {
		t.InputSchema.Properties[property.Name] = property
	}
	return t
}

func (t *Tool) invoke(arguments map[string]any) (string, error) {
	return t.handler(arguments)
}

func NewToolProperty(name, typ, description string) Property {
	return Property{
		Name:        name,
		Type:        typ,
		Description: description,
	}
}
