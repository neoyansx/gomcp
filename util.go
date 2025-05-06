package main

import (
	"fmt"
	"gomcp/protocol"
	"io"
	"log"
	"net/http"
)

func writeToClient(encoder protocol.Encoder, w io.Writer) error {
	result, err := encoder.MarshalJson()
	if err != nil {
		return err
	}
	st := fmt.Sprintf("event:message\ndata:%s\r\n\r\n", result)
	_, err = fmt.Fprint(w, st)
	if err != nil {
		log.Printf("error writing response: %v", err)
		return err
	}
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush() // for SSE write back immediately
	}
	return nil
}
