package protocol

import (
	"errors"
	"fmt"
	"net/url"
)

type (
	Resource struct {
		Uri          string `json:"uri"`
		Description  string `json:"description,omitempty"`
		Name         string `json:"name"`
		MimeType     string `json:"mimeType,omitempty"`
		handler      ResourceReader
		resourceType ResourceType
	}

	ResourceTemplate struct {
		UriTemplate string `json:"uriTemplate"`
		Description string `json:"description,omitempty"`
		Name        string `json:"name"`
		MimeType    string `json:"mimeType,omitempty"`
	}
)

func (r *Resource) getType() ResourceType {
	return r.resourceType
}

func (r *Resource) getMime() string {
	return r.MimeType
}

func (r *Resource) Mime(mime string) IResource {
	r.MimeType = mime
	return r
}

func (r *Resource) Content(typ ResourceType) IResource {
	r.resourceType = typ
	return r
}

func (r *Resource) Describe(description string) IResource {
	r.Description = description
	return r
}

func (r *Resource) Handler(handler ResourceReader) IResource {
	r.handler = handler
	return r
}

func (r *Resource) SetDescribe(description string) {
	r.Description = description
}

func (r *Resource) SetHandler(handler ResourceReader) {
	r.handler = handler
}

func (r *Resource) URI(uri string) IResource {
	r.Uri = uri
	return r
}

func (r *Resource) getUri() string {
	return r.Uri
}

func (r *Resource) Read() (string, error) {
	if _, err := url.ParseRequestURI(r.Uri); err != nil {
		return "", errors.New(fmt.Sprintf("Invalid resource uri: %s", r.Uri))
	}

	if r.handler == nil {
		return "", errors.New("resource reader is nil")
	}
	return r.handler(r.Uri)
}
