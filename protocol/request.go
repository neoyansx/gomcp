package protocol

import "errors"

type (
	Request struct {
		JsonRPC string         `json:"jsonrpc"`
		ID      any            `json:"id"`
		Method  string         `json:"method"`
		Params  map[string]any `json:"params,omitempty"`
	}
)

func (r Request) GetParam(param string) (any, error) {
	if r.Params == nil {
		return nil, errors.New("no params")
	}
	if v, ok := r.Params[param]; !ok {
		return nil, errors.New("param not found")
	} else {
		return v, nil
	}
}
