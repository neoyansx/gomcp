package protocol

import "errors"

func (s *Server) Prompt(name string) IPrompt {
	return &Prompt{
		Name: name,
	}
}

func (s *Server) AddPrompt(prompt IPrompt) {
	s.addPrompt(prompt)
}

func (s *Server) addPrompt(prompt IPrompt) {
	s.prompts.add(prompt)
}

func (s *Server) AddPrompts(prompts ...IPrompt) {
	s.addPrompts(prompts...)
}

func (s *Server) addPrompts(prompts ...IPrompt) {
	if len(prompts) == 0 {
		return
	}
	for _, prompt := range prompts {
		s.addPrompt(prompt)
	}
}

func (s *Server) GetPrompt(name string) (IPrompt, error) {
	return s.prompts.fetch(name)
}

func (s *Server) GetPrompts() ([]IPrompt, error) {
	if !s.listChangedPrompts {
		return nil, errors.New("does not support list prompts changes")
	}
	return s.prompts.fetchAll()
}
