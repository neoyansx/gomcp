package types

type (
	ClientParams struct {
		ProtocolVersion string              `json:"protocolVersion"`
		Capabilities    *ClientCapabilities `json:"capabilities"`
		ClientInfo      *ClientInfo         `json:"clientInfo"`
	}
	ClientCapabilities struct {
		Roots        *RootCapability `json:"roots"`
		Sampling     map[string]any  `json:"sampling"`
		Experimental map[string]any  `json:"experimental"`
	}

	RootCapability struct {
		ListChanged bool `json:"listChanged"`
	}

	ClientInfo struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
)
