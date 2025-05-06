package protocol

import "sync"

type (
	contents     []any
	ContentType  uint
	ResourceType uint
)

// for Tool and Prompt invoke result
const (
	Text ContentType = iota
	Image
	Audio
)

const (
	TextRes ResourceType = iota
	BinaryRes
)

type contentTypes struct {
	sync.RWMutex
	ts map[ContentType]string
}

func (c *contentTypes) add(typ ContentType, content string) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.ts[typ]; ok {
		return
	}
	c.ts[typ] = content
}

func (c *contentTypes) get(typ ContentType) string {
	c.RLock()
	defer c.RUnlock()
	content, ok := c.ts[typ]
	if !ok {
		return ""
	}
	return content
}

var ContentTypes = &contentTypes{
	ts: map[ContentType]string{
		Text:  "text",
		Image: "image",
		Audio: "audio",
	},
}

type (
	Content struct {
		Type     string    `json:"type"`
		Text     string    `json:"text,omitempty"`
		Data     string    `json:"data,omitempty"`
		MimeType string    `json:"mimeType,omitempty"`
		Resource *Resource `json:"resource,omitempty"`
	}
	ResourceContent struct {
		Uri      string `json:"uri"`
		Text     string `json:"text,omitempty"`
		Blob     string `json:"blob,omitempty"`
		MimeType string `json:"mimeType,omitempty"`
	}
)

func newResourceContent() IResourceContent {
	return &ResourceContent{}
}

func (r *ResourceContent) SetData(resType ResourceType, data string) {
	switch resType {
	case TextRes:
		r.Text = data
	case BinaryRes:
		r.Blob = data
	}
}

func (r *ResourceContent) SetMime(mime string) {
	r.MimeType = mime
}

func (r *ResourceContent) SetUri(uri string) {
	r.Uri = uri
}

func NewContent(typ ContentType) IContent {
	return newContent(typ)
}

func newContent(typ ContentType) IContent {
	t := ContentTypes.get(typ)
	if t == "" {
		return nil
	}
	return &Content{
		Type: t,
	}
}

func (c *Content) SetType(typ string) {
	c.Type = typ
}

func (c *Content) SetData(data string) {
	switch c.Type {
	case "image", "audio":
		c.Data = data
	case "text":
		c.Text = data
	}
}

func (c *Content) SetMime(mime string) {
	if c.Type == "text" {
		return
	}
	c.MimeType = mime
}

func (c *Content) SetResource(resource *Resource) {
	if len(c.Text) != 0 || len(c.Data) != 0 {
		return
	}
	c.Resource = resource
}

func (c *contents) add(content any) {
	*c = append(*c, content)
}
