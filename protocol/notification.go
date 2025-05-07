package protocol

import (
	"encoding/json"
	"github.com/neoyansx/gomcp/common"
	"log"
	"sync"
)

type Notification struct {
	JsonRPC string         `json:"jsonrpc"`
	Method  string         `json:"method"`
	Meta    map[string]any `json:"_meta,omitempty"`
	Params  map[string]any `json:"params,omitempty"`
}

func NewNotification(method string) *Notification {
	return &Notification{
		JsonRPC: common.JsonRPCVersion,
		Method:  method,
	}
}

func (n *Notification) Write(key string, value any) {
	mu := &sync.RWMutex{}
	mu.Lock()
	defer mu.Unlock()
	n.Params[key] = value
}

func (n *Notification) MarshalJson() ([]byte, error) {
	data, err := json.Marshal(n)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return data, nil
}

func (n *Notification) String() string {
	data, _ := n.MarshalJson()
	return string(data)
}
