package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"gomcp/protocol"
	"log"
	"net/http"
	"net/url"
	"time"
)

type (
	SSE struct {
		url         url.URL
		accessPoint string
		mcpSrv      protocol.IServer
		chResult    chan *protocol.Response
		sessions    *sessions
		httpHandler http.Handler
	}
)

func NewSSEServer(url, accessPoint url.URL, mcpSrv protocol.IServer) *SSE {
	return &SSE{
		url:         url,
		accessPoint: accessPoint.Path,
		chResult:    make(chan *protocol.Response),
		sessions: &sessions{
			se: make(map[string]*session),
		},
		mcpSrv: mcpSrv,
	}
}
func (s *SSE) Router(router http.Handler) {
	s.httpHandler = router
}

func (s *SSE) router() *chi.Mux {
	routers := chi.NewRouter()
	routers.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Content-Disposition", "Content-Length"},
		AllowCredentials: true,
	}))
	routers.Use(middleware.Logger)
	routers.Use(middleware.Recoverer)
	routers.Use(middleware.RedirectSlashes)
	routers.Use(middleware.StripSlashes)
	routers.Use(middleware.Timeout(60 * time.Minute))
	routers.Get("/", s.Listen())
	routers.Post(s.accessPoint, s.Messages())
	return routers
}

func (s *SSE) Run() error {
	log.Printf("listening on %s", s.url.Port())
	if s.httpHandler != nil {
		return http.ListenAndServe(fmt.Sprintf(":%s", s.url.Port()), s.httpHandler)

	}
	return http.ListenAndServe(fmt.Sprintf(":%s", s.url.Port()), s.router())
}
