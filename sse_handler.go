package gomcp

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/neoyansx/gomcp/protocol"
	"io"
	"log"
	"net/http"
	"strings"
)

func (s *SSE) Listen() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("cache-control", "no-cache,no-transform")
		w.Header().Set("content-type", "text/event-stream")
		w.Header().Set("connection", "keep-alive")
		w.WriteHeader(http.StatusOK)

		message := make(chan *protocol.Response)
		sessionID := strings.ReplaceAll(uuid.New().String(), "-", "")
		session := &session{
			done:         make(chan struct{}),
			notification: make(chan *protocol.Notification),
			resp:         message,
		}

		s.sessions.add(sessionID, session)

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}
		st := fmt.Sprintf("event:endpoint\ndata:%s\r\n\r\n", fmt.Sprintf("%s?sessionId=%s", s.accessPoint, sessionID))
		_, err := fmt.Fprint(w, st)
		if err != nil {
			log.Printf("error writing response: %v", err)
			return
		}
		flusher.Flush()

		go func(sessionID string) {
			select {
			case <-r.Context().Done():
				s.sessions.remove(sessionID)
				return
			case <-session.done:
				return
			}
		}(sessionID)

		for {
			select {
			case <-r.Context().Done():
				close(session.done)
				return
			case resp, ok := <-session.resp:
				if !ok {
					return
				}
				err := writeToClient(resp, w)
				if err != nil {
					return
				}
			case notification, ok := <-session.notification:
				if !ok {
					return
				}
				err := writeToClient(notification, w)
				if err != nil {
					return
				}
			}
		}
	}
}

func (s *SSE) Messages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.URL.Query().Get("sessionId")
		if sessionID == "" {
			log.Printf("sessionId not found in query")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		session, err := s.sessions.get(sessionID)
		if err != nil {
			log.Printf("error getting session: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		defer r.Body.Close()
		log.Printf("body: %s", body)
		req := protocol.Request{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Println(err)
			return
		}
		if req.JsonRPC != "2.0" {
			log.Println("expected 2.0, got ", req.JsonRPC)
			return
		}
		/*resp := s.mcpSrv.ProcessRequest(req)
		session.resp <- resp*/
		session.resp <- s.mcpSrv.InvokeHandler(req)
		w.WriteHeader(http.StatusAccepted)
	}
}

func (s *SSE) SendNotification(notification *protocol.Notification) {
	s.sendNotification(notification)
}

func (s *SSE) sendNotification(notification *protocol.Notification) {
	sessions, err := s.sessions.list()
	if err != nil {
		return
	}
	for _, session := range sessions {
		session.notification <- notification
	}
}
