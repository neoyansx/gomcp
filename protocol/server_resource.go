package protocol

import "errors"

type ResourceReader func(uri string) (string, error)

func (s *Server) Resource(name string) IResource {
	return &Resource{
		Name: name,
	}
}

func (s *Server) removeResource(resource IResource) error {
	return s.resources.remove(resource.getUri())
}

func (s *Server) AddResource(resource IResource) {
	s.addResource(resource)
}

func (s *Server) addResource(resource IResource) {
	uri := resource.getUri()
	if s.resources.exists(uri) {
		return
	}
	s.resources.add(uri, resource)
}

func (s *Server) AddResources(resources ...IResource) {
	s.addResources(resources...)
}

func (s *Server) addResources(resources ...IResource) {
	if len(resources) == 0 {
		return
	}
	for _, resource := range resources {
		if s.resources.exists(resource.getUri()) {
			continue
		}
		s.resources.add(resource.getUri(), resource)
	}
}

func (s *Server) GetResources() ([]IResource, error) {
	if !s.isListChangedResources() {
		return nil, errors.New("does not support list changed resources")
	}
	return s.resources.fetchAll()
}

func (s *Server) GetResource(uri string) (IResource, error) {
	return s.resources.fetch(uri)
}
