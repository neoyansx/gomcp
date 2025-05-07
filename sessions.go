package main

import (
	"errors"
	"github.com/neoyansx/gomcp/protocol"
	"sync"
)

type (
	session struct {
		reqID        any
		done         chan struct{}
		resp         chan *protocol.Response
		notification chan *protocol.Notification
	}
	sessions struct {
		sync.RWMutex
		se map[string]*session
	}
)

func newSessions() *sessions {
	return &sessions{
		se: make(map[string]*session),
	}
}

func (s *sessions) add(sessionID string, session *session) {
	s.Lock()
	defer s.Unlock()
	s.se[sessionID] = session
}

func (s *sessions) get(sessionID string) (*session, error) {
	s.RLock()
	defer s.RUnlock()
	if resp, ok := s.se[sessionID]; !ok {
		return nil, errors.New("session not found")
	} else {
		return resp, nil
	}
}

func (s *sessions) remove(sessionID string) error {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.se[sessionID]; ok {
		delete(s.se, sessionID)
		return nil
	} else {
		return errors.New("session id not found")
	}
}

func (s *sessions) resetSession(sessionID string, session *session) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.se[sessionID]; ok {
		s.se[sessionID] = session
	}
}

func (s *sessions) exists(sessionID string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.se[sessionID]
	return ok
}

func (s *sessions) list() ([]*session, error) {
	s.RLock()
	defer s.RUnlock()
	if s.se == nil {
		return nil, errors.New("sessions is empty")
	}
	list := make([]*session, 0, len(s.se))
	for _, session := range s.se {
		list = append(list, session)
	}
	return list, nil
}

func (s *sessions) getSessionByRequestID(reqID any) (*session, error) {
	ss, err := s.list()
	if err != nil {
		return nil, err
	}
	var session *session
	for _, item := range ss {
		if item.reqID == reqID {
			session = item
			break
		}
	}
	return session, nil
}
