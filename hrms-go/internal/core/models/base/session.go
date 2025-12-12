package base

import (
	"time"

	"gorm.io/gorm"
)

// Session represents a user session (for both cookie and JWT-based auth)
type Session struct {
	gorm.Model
	UserID         uint      `gorm:"not null;index" json:"user_id"`
	User           User      `gorm:"foreignKey:UserID" json:"user,omitempty"`

	// Session identification
	SessionID      string    `gorm:"uniqueIndex;not null;size:255" json:"session_id"`
	RefreshToken   string    `gorm:"uniqueIndex;size:500" json:"-"`

	// Session metadata
	IPAddress      string    `gorm:"size:45" json:"ip_address"`
	UserAgent      string    `gorm:"type:text" json:"user_agent"`
	DeviceType     string    `gorm:"size:50" json:"device_type"` // mobile, desktop, tablet
	DeviceName     string    `gorm:"size:100" json:"device_name"`

	// Session timing
	LastActivity   time.Time `gorm:"index" json:"last_activity"`
	ExpiresAt      time.Time `gorm:"index" json:"expires_at"`

	// Session status
	Active         bool      `gorm:"default:true;index" json:"active"`
	LoggedOut      bool      `gorm:"default:false" json:"logged_out"`
	LoggedOutAt    *time.Time `json:"logged_out_at,omitempty"`

	// Security
	IsRevoked      bool      `gorm:"default:false" json:"is_revoked"`
	RevokedAt      *time.Time `json:"revoked_at,omitempty"`
	RevokeReason   string    `gorm:"type:text" json:"revoke_reason,omitempty"`
}

// IsExpired checks if the session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsValid checks if session is valid and active
func (s *Session) IsValid() bool {
	return s.Active && !s.LoggedOut && !s.IsRevoked && !s.IsExpired()
}

// Revoke revokes the session
func (s *Session) Revoke(reason string) {
	s.IsRevoked = true
	now := time.Now()
	s.RevokedAt = &now
	s.RevokeReason = reason
	s.Active = false
}

// Logout marks the session as logged out
func (s *Session) Logout() {
	s.LoggedOut = true
	now := time.Now()
	s.LoggedOutAt = &now
	s.Active = false
}

// UpdateActivity updates the last activity timestamp
func (s *Session) UpdateActivity() {
	s.LastActivity = time.Now()
}

// TableName specifies the table name for Session model
func (Session) TableName() string {
	return "sessions"
}

// RefreshTokenClaims represents claims stored in refresh token
type RefreshTokenClaims struct {
	UserID    uint   `json:"user_id"`
	SessionID string `json:"session_id"`
	Email     string `json:"email"`
}
