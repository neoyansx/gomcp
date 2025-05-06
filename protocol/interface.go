package protocol

type (
	IServer interface {
		AddHandler(method string, hd handler)
		SetHandler(method string, hd handler)
		RemoveHandler(method string)
		Tool(name string) ITool
		AddTools(tools ...ITool)
		Prompt(name string) IPrompt
		AddPrompts(prompts ...IPrompt)
		Resource(name string) IResource
		AddResources(resources ...IResource)
		GetTool(name string) (ITool, error)
		GetTools() ([]ITool, error)
		GetPrompt(name string) (IPrompt, error)
		GetResource(name string) (IResource, error)
		InvokeHandler(req Request) *Response
	}

	Descriptor interface {
		SetDescribe(description string)
	}

	HandlerSetter interface {
		SetHandler(handler HandleFunc)
	}

	ContentGetter interface {
		getContent() ContentType
	}

	ITool interface {
		Descriptor
		HandlerSetter
		ContentGetter
		Describe(description string) ITool
		Handler(handler HandleFunc) ITool
		Property(property Property) ITool
		Properties(ps ...Property) ITool
		Content(typ ContentType) ITool
		invoke(arguments map[string]any) (string, error)
		getMime() string
		getName() string
	}

	IPrompt interface {
		Descriptor
		HandlerSetter
		ContentGetter
		Describe(description string) IPrompt
		Handler(handler HandleFunc) IPrompt
		Role(role string) IPrompt
		Content(typ ContentType) IPrompt
		Argument(name string, required bool) IPrompt
		getName() string
		getMime() string
		invoke(arguments map[string]any) (string, error)
		getRole() string
		getDescription() string
	}

	IResource interface {
		Descriptor
		Describe(description string) IResource
		URI(uri string) IResource
		Handler(handler ResourceReader) IResource
		Mime(mime string) IResource
		SetHandler(handler ResourceReader)
		Content(typ ResourceType) IResource
		Read() (string, error)
		getType() ResourceType
		getMime() string
		getUri() string
	}
	// IContent interface for Tool and Prompt content
	IContent interface {
		SetType(typ string)
		// SetData set data or text field
		SetData(data string)
		SetMime(mime string)
		SetResource(resource *Resource)
	}
	// IResourceContent interface for Resource content
	IResourceContent interface {
		// SetData set text or blog field
		SetData(resType ResourceType, data string)
		SetMime(mime string)
		SetUri(uri string)
	}

	Encoder interface {
		MarshalJson() ([]byte, error)
	}
	// Writer write data to client the param is a key-value pair
	Writer interface {
		Write(key string, value any)
	}
)
