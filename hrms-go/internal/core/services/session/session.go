package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
)

// Session represents a user session
type Session struct {
	SID        string    `json:"sid"`
	UserID     string    `json:"user_id"`
	User       string    `json:"user"`
	FullName   string    `json:"full_name"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	LastActive time.Time `json:"last_active"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
}

// Store interface for session storage
type Store interface {
	Create(ctx context.Context, session *Session) error
	Get(ctx context.Context, sid string) (*Session, error)
	Update(ctx context.Context, session *Session) error
	Delete(ctx context.Context, sid string) error
	DeleteByUser(ctx context.Context, userID string) error
}

// MemoryStore implements in-memory session storage
type MemoryStore struct {
	sessions map[string]*Session
	mu       sync.RWMutex
}

// NewMemoryStore creates a new in-memory session store
func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		sessions: make(map[string]*Session),
	}

	// Start cleanup goroutine
	go store.cleanup()

	return store
}

// Create stores a new session
func (s *MemoryStore) Create(ctx context.Context, session *Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[session.SID] = session
	return nil
}

// Get retrieves a session by SID
func (s *MemoryStore) Get(ctx context.Context, sid string) (*Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[sid]
	if !exists {
		return nil, ErrSessionNotFound
	}

	// Check if expired
	if time.Now().After(session.ExpiresAt) {
		return nil, ErrSessionExpired
	}

	return session, nil
}

// Update updates an existing session
func (s *MemoryStore) Update(ctx context.Context, session *Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.sessions[session.SID]; !exists {
		return ErrSessionNotFound
	}

	s.sessions[session.SID] = session
	return nil
}

// Delete removes a session
func (s *MemoryStore) Delete(ctx context.Context, sid string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, sid)
	return nil
}

// DeleteByUser removes all sessions for a user
func (s *MemoryStore) DeleteByUser(ctx context.Context, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for sid, session := range s.sessions {
		if session.UserID == userID {
			delete(s.sessions, sid)
		}
	}

	return nil
}

// cleanup removes expired sessions periodically
func (s *MemoryStore) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for sid, session := range s.sessions {
			if now.After(session.ExpiresAt) {
				delete(s.sessions, sid)
			}
		}
		s.mu.Unlock()
	}
}

// Service provides session management functionality
type Service struct {
	store          Store
	sessionTimeout time.Duration
}

// NewService creates a new session service
func NewService(store Store, sessionTimeout time.Duration) *Service {
	if sessionTimeout == 0 {
		sessionTimeout = 24 * time.Hour // Default 24 hours
	}

	return &Service{
		store:          store,
		sessionTimeout: sessionTimeout,
	}
}

// Create creates a new session for a user
func (s *Service) Create(ctx context.Context, userID, username, fullName, email, ipAddress, userAgent string) (*Session, error) {
	sid, err := generateSID()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	session := &Session{
		SID:        sid,
		UserID:     userID,
		User:       username,
		FullName:   fullName,
		Email:      email,
		CreatedAt:  now,
		ExpiresAt:  now.Add(s.sessionTimeout),
		LastActive: now,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
	}

	err = s.store.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// Get retrieves a session and updates last active time
func (s *Service) Get(ctx context.Context, sid string) (*Session, error) {
	session, err := s.store.Get(ctx, sid)
	if err != nil {
		return nil, err
	}

	// Update last active time
	session.LastActive = time.Now()
	_ = s.store.Update(ctx, session)

	return session, nil
}

// Delete removes a session (logout)
func (s *Service) Delete(ctx context.Context, sid string) error {
	return s.store.Delete(ctx, sid)
}

// DeleteAllForUser removes all sessions for a user
func (s *Service) DeleteAllForUser(ctx context.Context, userID string) error {
	return s.store.DeleteByUser(ctx, userID)
}

// Extend extends a session's expiry time
func (s *Service) Extend(ctx context.Context, sid string) error {
	session, err := s.store.Get(ctx, sid)
	if err != nil {
		return err
	}

	session.ExpiresAt = time.Now().Add(s.sessionTimeout)
	return s.store.Update(ctx, session)
}

// generateSID generates a random session ID
func generateSID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
