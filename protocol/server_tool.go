package protocol

import "errors"

func (s *Server) Tool(name string) ITool {
	return &Tool{
		Name:        name,
		InputSchema: newInputSchema(),
	}
}

func (s *Server) AddTool(tool ITool) {
	s.tools.add(tool)
}

func (s *Server) AddTools(tools ...ITool) {
	s.tools.addMore(tools...)
}

func (s *Server) GetTool(name string) (ITool, error) {
	return s.tools.fetch(name)
}

func (s *Server) RemoveTool(name string) error {
	return s.tools.remove(name)
}

func (s *Server) GetTools() ([]ITool, error) {
	if !s.listChangedTools {
		return nil, errors.New("does not support list changed tools")
	}
	return s.tools.fetchAll()
}
