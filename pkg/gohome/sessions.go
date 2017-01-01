package gohome

import (
	"encoding/base64"
	"math/rand"
	"sync"
)

// Sessions manages user sessions in the app
type Sessions struct {
	mutex sync.RWMutex
	sids  map[string]map[string]interface{}
}

// NewSessions returns a newly instantiated Sessions instance
func NewSessions() *Sessions {
	return &Sessions{
		sids: make(map[string]map[string]interface{}),
	}
}

// Add adds a new session and returns the session ID back to the caller.  You must
// call Save() at some point to persist the sessions to disk
func (s *Sessions) Add() (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	b := make([]byte, 128)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	//TODO: we need to have a expiration date in here too
	sid := base64.URLEncoding.EncodeToString(b)
	s.sids[sid] = make(map[string]interface{})
	return sid, nil
}

// Get returns the values associated with the session ID.  If the session ID is
// not valid it returns false as the second return value
func (s *Sessions) Get(SID string) (map[string]interface{}, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	val, ok := s.sids[SID]

	return val, ok
}

// Save persists the session information to backing store
func (s *Sessions) Save() error {
	//TODO: In memory right now, need to persist to disk otherwise sessions are lost
	//across reboots
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return nil
}
