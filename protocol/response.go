package protocol

import (
	"encoding/json"
	"errors"
	"fmt"
	"gomcp/common"
	"log"
	"sync"
)

type (
	Response struct {
		JsonRPC string         `json:"jsonrpc"`
		ID      any            `json:"id"`
		Result  map[string]any `json:"result"`
		Error   *Error         `json:"error,omitempty"`
	}

	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    any    `json:"data,omitempty"`
	}
)

func NewResponse(id any) *Response {
	switch id.(type) {
	case string, int, float64, float32:
		return &Response{
			JsonRPC: common.JsonRPCVersion,
			ID:      id,
			Result:  map[string]any{"protocolVersion": common.ProtocolVersion},
		}
	default:
		return nil
	}
}

func (r *Response) Write(key string, value any) {
	mu := &sync.RWMutex{}
	mu.Lock()
	defer mu.Unlock()
	r.Result[key] = value
}

func (r *Response) MarshalJson() ([]byte, error) {
	data, err := json.Marshal(r)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return data, nil
}

func (r *Response) WriteError(code int, err error) error {
	if r.Result != nil {
		return errors.New("cannot write error because result was set")
	}
	e := &Error{
		Code:    code,
		Message: err.Error(),
	}
	r.Error = e
	return nil
}

func (r *Response) String() string {
	data, _ := r.MarshalJson()
	return string(data)
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error %d: %s data: %v", e.Code, e.Message, e.Data)
}
