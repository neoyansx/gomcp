package protocol

import "log"

type (
	// HandleFunc provide external function calls for Tool and Prompt only
	HandleFunc func(map[string]any) (string, error)
)

const (
	ServerInit    = "initialize"
	ListPrompts   = "prompts/list"
	GetPrompt     = "prompts/get"
	ListTools     = "tools/list"
	InvokeTool    = "tools/call"
	ToolsChanged  = "notifications/tools/list_changed"
	ListResources = "resources/list"
	GetResource   = "resources/read"
)

func (s *Server) Initialize(req Request) *Response {
	return s.initialize(req)
}

func (s *Server) initialize(req Request) *Response {
	cas := Capabilities{
		LoggingCapability: &LoggingCapability{},
		PromptCapability: &PromptCapability{
			ListChanged: s.listChangedPrompts,
		},
		ToolsCapability: &ToolsCapability{
			ListChanged: s.listChangedTools,
		},
		ResourceCapability: &ResourceCapability{
			Subscribe:   s.subscribeResource,
			ListChanged: s.listChangedResource,
		},
	}
	response := NewResponse(req.ID)
	response.Write("capabilities", cas)
	response.Write("serverInfo", s.serverInfo)
	return response
}

func (s *Server) FetchPrompts(req Request) *Response {
	return s.fetchPrompts(req)
}

func (s *Server) fetchPrompts(req Request) *Response {
	resp := NewResponse(req.ID)
	prompts, err := s.GetPrompts()
	if err != nil {
		log.Printf("failed to list prompts: %v", err)
		resp.WriteError(INTERNAL_ERROR, err)
	} else {
		resp.Write("prompts", prompts)
	}
	return resp
}

func (s *Server) FetchTools(req Request) *Response {
	return s.fetchTools(req)
}

func (s *Server) fetchTools(req Request) *Response {
	resp := NewResponse(req.ID)
	tools, err := s.GetTools()
	if err != nil {
		resp.WriteError(INTERNAL_ERROR, err)
		return resp
	}
	resp.Write("tools", tools)
	return resp
}

func (s *Server) FetchResources(req Request) *Response {
	return s.fetchResources(req)
}

func (s *Server) fetchResources(req Request) *Response {
	resp := NewResponse(req.ID)
	resources, err := s.GetResources()
	if err != nil {
		log.Printf("failed to list resources: %v", err)
		resp.WriteError(INTERNAL_ERROR, err)
	} else {
		resp.Write("resources", resources)
	}
	return resp
}

func (s *Server) FetchPrompt(req Request) *Response {
	return s.fetchPrompt(req)
}

func (s *Server) fetchPrompt(req Request) *Response {
	var (
		resp      = NewResponse(req.ID)
		name, err = req.GetParam("name")
	)
	if err != nil {
		resp.WriteError(INVALID_REQUEST, err)
		return resp
	}
	prompt, err := s.GetPrompt(name.(string))
	if err != nil {
		resp.WriteError(INVALID_REQUEST, err)
		return resp
	}
	args, err := req.GetParam("arguments")
	if err != nil {
		resp.WriteError(INVALID_REQUEST, err)
		return resp
	}
	messages := make(messages, 0)

	result, err := prompt.invoke(args.(map[string]any))
	if err != nil {
		resp.WriteError(INVALID_REQUEST, err)
		return resp
	}
	resp.Write("description", prompt.getDescription())
	var (
		msg         = NewPromptMessage(prompt.getRole())
		contentType = prompt.getContent()
		content     = newContent(contentType)
	)
	content.SetData(result)
	switch contentType {
	case Image, Audio:
		content.SetMime(prompt.getMime())
		msg.AddContent(content)
	case Text:
		msg.AddContent(content)
	}
	messages.addMessage(msg)
	resp.Write("messages", messages)
	return resp
}

func (s *Server) FetchResource(req Request) *Response {
	return s.fetchResource(req)
}

func (s *Server) fetchResource(req Request) *Response {
	var (
		resp     = NewResponse(req.ID)
		uri, err = req.GetParam("uri")
	)
	if err != nil {
		log.Println(err)
		resp.WriteError(INTERNAL_ERROR, err)
		return resp
	}

	resource, err := s.GetResource(uri.(string))
	if err != nil {
		log.Println(err)
		resp.WriteError(INVALID_REQUEST, err)
		return resp
	}
	result, err := resource.Read()
	if err != nil {
		log.Println(err)
		resp.WriteError(INTERNAL_ERROR, err)
		return resp
	}
	var (
		resContents = make(contents, 0)
		resType     = resource.getType()
		content     = newResourceContent()
	)
	content.SetData(resType, result)
	content.SetMime(resource.getMime())
	resContents = append(resContents, content)
	resp.Write("contents", resContents)
	return resp
}

func (s *Server) InvokeTool(req Request) *Response {
	return s.invokeTool(req)
}

func (s *Server) invokeTool(req Request) *Response {
	var (
		resp     = NewResponse(req.ID)
		contents = make(contents, 0)
	)
	args, err := req.GetParam("arguments")
	if err != nil {
		log.Println(err)
		resp.WriteError(INVALID_REQUEST, err)
		return resp
	}
	log.Printf("args: %v", args)
	arguments := args.(map[string]any)
	toolName, err := req.GetParam("name")
	if err != nil {
		log.Println(err)
		resp.WriteError(INVALID_REQUEST, err)
		return resp
	}
	tool, err := s.tools.fetch(toolName.(string))
	if err != nil {
		log.Println(err)
		resp.WriteError(INTERNAL_ERROR, err)
		return resp
	}
	result, err := tool.invoke(arguments)
	if err != nil {
		log.Println(err)
		resp.WriteError(INTERNAL_ERROR, err)
		return resp
	}
	contentType := tool.getContent()
	content := newContent(contentType)
	content.SetData(result)
	if contentType == Image || contentType == Audio {
		content.SetMime(tool.getMime())
	}
	contents.add(content)
	resp.Write("contents", contents)
	return resp
}
