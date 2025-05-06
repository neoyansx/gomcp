package protocol

const (
	PARSE_ERROR      = -32700
	INVALID_REQUEST  = -32600
	METHOD_NOT_FOUND = -32601
	INVALID_PARAMS   = -32602
	INTERNAL_ERROR   = -32603
)

type ComponentType uint

const (
	CmpTool ComponentType = iota
	CmpPrompt
	CmpResource
)

var _IServer = (*Server)(nil)

type (
	Server struct {
		tools               *tools
		resources           *resources
		serverInfo          *ServerInfo
		prompts             *prompts
		listChangedTools    bool
		listChangedPrompts  bool
		listChangedResource bool
		subscribeResource   bool
		chNotifications     chan *Notification
		hds                 *handlers
	}

	MCPServerInitResult struct {
		ProtocolVersion string        `json:"protocolVersion"`
		Capabilities    *Capabilities `json:"capabilities"`
		ServerInfo      *ServerInfo   `json:"serverInfo"`
	}

	Capabilities struct {
		*LoggingCapability  `json:"logging"`
		*PromptCapability   `json:"prompts"`
		*ToolsCapability    `json:"tools"`
		*ResourceCapability `json:"resources"`
	}

	ServerInfo struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	LoggingCapability struct {
		Logging map[string]any `json:"-"`
	}

	PromptCapability struct {
		ListChanged bool `json:"listChanged"`
	}

	ToolsCapability struct {
		ListChanged bool `json:"listChanged"`
	}

	ResourceCapability struct {
		Subscribe   bool `json:"subscribe,omitempty"`
		ListChanged bool `json:"listChanged,omitempty"`
	}
)

func NewServer(name, version string) IServer {
	srv := &Server{
		serverInfo: &ServerInfo{
			Name:    name,
			Version: version,
		},
		tools: &tools{
			ts: make(map[string]ITool),
		},
		resources: &resources{
			rs: make(map[string]IResource),
		},
		prompts: &prompts{
			ps: make(map[string]IPrompt),
		},
		chNotifications: make(chan *Notification),
		hds: &handlers{
			hs: make(map[string]handler),
		},
	}
	srv.AddHandler(ServerInit, srv.initialize)
	srv.AddHandler(ListPrompts, srv.fetchPrompts)
	srv.AddHandler(ListTools, srv.fetchTools)
	srv.AddHandler(ListResources, srv.fetchResources)
	srv.AddHandler(GetPrompt, srv.fetchPrompt)
	srv.AddHandler(GetResource, srv.fetchResource)
	srv.AddHandler(InvokeTool, srv.invokeTool)
	return srv
}

func (s *Server) InvokeHandler(req Request) *Response {
	return s.hds.invoke(req)
}

func (s *Server) AddHandler(method string, hd handler) {
	s.hds.add(method, hd)
}

func (s *Server) SetHandler(method string, hd handler) {
	s.hds.set(method, hd)
}

func (s *Server) RemoveHandler(method string) {
	s.hds.remove(method)
}

// AddContentTypes add customized  content type for Tool and Prompt invoke
func (s *Server) AddContentTypes(typ ContentType, name string) {
	ContentTypes.add(typ, name)
}

func (s *Server) Capabilities(listChangedTools, listChangedPrompts, listChangedSources bool, subscribeResources bool) {
	s.listChangedTools = listChangedTools
	s.listChangedPrompts = listChangedSources
	s.listChangedResource = listChangedSources
	s.subscribeResource = subscribeResources
}

func (s *Server) isListChangedTools() bool {
	return s.listChangedTools
}

func (s *Server) isListChangedResources() bool {
	return s.listChangedResource
}

func (s *Server) isListChangedPrompts() bool {
	return s.listChangedPrompts
}

func (s *Server) isSubscribeResource() bool {
	return s.subscribeResource
}

func (s *Server) ProcessRequest(req Request) *Response {
	return s.hds.invoke(req)
	/*switch req.Method {
	case initialize:
		return s.initialize(req.ID)
	case "notifications/initialized":
		log.Println("notifications/initialized")
		return &Response{}

	case "ping":
		return &Response{}

	case "completion/complete":
		log.Println("completion/complete")
		return nil

	case "logging/setLevel":
		log.Println("logging/setLevel")
		return nil

	case "notifications/cancelled":
		log.Println("notifications/cancelled")
		resp := &Response{}
		return resp
	case "notifications/message":
		log.Println("notifications/message")
		return nil

	case "notifications/progress":
		log.Println("notifications/progress")
		return nil

	case "notifications/prompts/list_changed":
		log.Println("notifications/prompts/list_changed")
		return nil

	case "notifications/resources/list_changed":
		log.Println("notifications/resources/list_changed")
		return nil

	case "notifications/resources/updated":
		return nil

	case "notifications/roots/list_changed":
		return nil

	case "notifications/tools/list_changed":
		return nil

	case "prompts/get":
		var (
			resp      = NewResponse(req.ID)
			name, err = req.GetParam("name")
		)
		if err != nil {
			resp.WriteError(INVALID_REQUEST, err)
			return resp
		}
		prompt, err := s.getPrompt(name.(string))
		if err != nil {
			resp.WriteError(INVALID_REQUEST, err)
			return resp
		}
		args, err := req.GetParam("arguments")
		if err != nil {
			e := Error{
				Code:    INVALID_REQUEST,
				Message: err.Error(),
			}
			resp.WriteError(&e)
			return resp
		}
		messages := make(messages, 0)

		result, err := prompt.generator(args.(map[string]any))
		if err != nil {
			e := Error{
				Code:    INVALID_REQUEST,
				Message: err.Error(),
			}
			resp.WriteError(&e)
			return resp
		}
		resp.Write("description", prompt.Description)
		msg := NewPromptMessage(prompt.role)
		switch {
		case prompt.isTextResult:
			msg.Text(result, prompt.mimeType)
			messages.addMessage(msg)
		case prompt.isAudioResult:
			msg.Audio(result, prompt.mimeType)
			messages.addMessage(msg)
		case prompt.isImageResult:
			msg.Image(result, prompt.mimeType)
			messages.addMessage(msg)
		}
		resp.Write("messages", messages)
		return resp
	case "prompts/list":
		resp := NewResponse(req.ID)
		prompts, err := s.ListPrompts()
		if err != nil {
			log.Printf("failed to list prompts: %v", err)
			e := Error{
				Code:    INTERNAL_ERROR,
				Message: err.Error(),
			}
			resp.WriteError(&e)
		} else {
			resp.Write("prompts", prompts)
		}
		return resp

	case "resources/list":
		resp := NewResponse(req.ID)
		resources, err := s.ListResources()
		if err != nil {
			log.Printf("failed to list resources: %v", err)
			e := Error{
				Code:    INTERNAL_ERROR,
				Message: err.Error(),
			}
			resp.WriteError(&e)
		} else {
			resp.Write("resources", resources)
		}
		return resp

	case "resources/read":
		var (
			resp     = NewResponse(req.ID)
			uri, err = req.GetParam("uri")
		)
		if err != nil {
			log.Println(err)
			e := Error{
				Code:    INTERNAL_ERROR,
				Message: err.Error(),
			}
			resp.WriteError(&e)
			return resp
		}

		resource, err := s.getResource(uri.(string))
		if err != nil {
			log.Println(err)
			e := Error{
				Code:    INVALID_REQUEST,
				Message: err.Error(),
			}
			resp.WriteError(&e)
			return resp
		}
		content, err := resource.Read()
		if err != nil {
			log.Println(err)
			e := Error{
				Code:    INTERNAL_ERROR,
				Message: err.Error(),
			}
			resp.WriteError(&e)
			return resp
		}
		resContents := make(contents, 0)
		switch {
		case resource.isText:
			text := ResTextContent{
				Uri:      resource.Uri,
				MimeType: resource.MIMEType,
				Text:     content,
			}
			resContents = append(resContents, text)
		case resource.isBinary:
			binary := ResBinaryContent{
				Uri:      resource.Uri,
				Blob:     content,
				MimeType: resource.MIMEType,
			}
			resContents = append(resContents, binary)
		default:
			e := Error{
				Code:    INVALID_REQUEST,
				Message: "unknown resource type",
			}
			resp.WriteError(&e)
			return resp
		}
		resp.Write("contents", resContents)
		return resp

	case "resources/subscribe":
		return nil

	case "resources/templates/list":
		return nil

	case "resources/unsubscribe":
		return nil

	case "roots/list":
		return nil

	case "sampling/createMessage":
		return nil

	case "tools/call":
		var (
			resp     = NewResponse(req.ID)
			contents = make(contents, 0)
		)
		args, err := req.GetParam("arguments")
		if err != nil {
			log.Println(err)
			e := Error{
				Code:    INVALID_REQUEST,
				Message: err.Error(),
			}
			resp.WriteError(&e)
			return resp
		}
		arguments := args.(map[string]any)
		toolName, err := req.GetParam("name")
		if err != nil {
			log.Println(err)
			e := Error{
				Code:    INVALID_REQUEST,
				Message: err.Error(),
			}
			resp.WriteError(&e)
			return resp
		}
		tool, err := s.tools.get(toolName.(string))
		if err != nil {
			log.Println(err)
			e := Error{
				Code:    INTERNAL_ERROR,
				Message: err.Error(),
			}
			resp.WriteError(&e)
			return resp
		}
		result, err := tool.do(arguments)
		if err != nil {
			log.Println(err)
			e := Error{
				Code:    INTERNAL_ERROR,
				Message: err.Error(),
			}
			resp.WriteError(&e)
			return resp
		}
		switch {
		case tool.isTextResult:
			txtResult := &TextResult{
				Type: "text",
				Text: result,
			}
			contents.add(txtResult)
			resp.Write("contents", txtResult)
		case tool.isAudioResult:
			audioResult := &MediaResult{
				Type: "audio",
				Data: result,
			}
			contents.add(audioResult)
			resp.Write("contents", audioResult)
		case tool.isImageResult:
			imageResult := &MediaResult{
				Type:     "image",
				Data:     result,
				MimeType: tool.resultType,
			}
			contents.add(imageResult)
			resp.Write("contents", imageResult)
		}

		return resp
	case "tools/list":
		tools, err := s.ListTools()
		if err != nil {
			log.Printf("failed to list tools: %v", err)
			return nil
		}
		response := NewResponse(req.ID)
		response.Write("tools", tools)
		return response

	default:
		return nil
	}*/
}
