package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/OtchereDev/hrms-go/internal/config"
	"github.com/OtchereDev/hrms-go/internal/core/models/base"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserDisabled       = errors.New("user account is disabled")
	ErrUserLocked         = errors.New("user account is locked")
	ErrUserNotFound       = errors.New("user not found")
	ErrSessionNotFound    = errors.New("session not found")
	ErrSessionExpired     = errors.New("session has expired")
	ErrSessionRevoked     = errors.New("session has been revoked")
)

// AuthService handles authentication operations
type AuthService struct {
	db         *gorm.DB
	jwtService *JWTService
	config     *config.Config
}

// NewAuthService creates a new authentication service
func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{
		db:         db,
		jwtService: NewJWTService(&cfg.JWT),
		config:     cfg,
	}
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"usr" validate:"required,email"` // Frappe uses "usr" field name
	Password string `json:"pwd" validate:"required"`        // Frappe uses "pwd" field name
}

// LoginResponse represents login response
type LoginResponse struct {
	User         *UserResponse `json:"user"`
	TokenPair    *TokenPair    `json:"token"`
	Message      string        `json:"message"`
	FullName     string        `json:"full_name"`
	Home_Page    string        `json:"home_page,omitempty"`
}

// UserResponse represents user data in response
type UserResponse struct {
	ID            uint       `json:"id"`
	Email         string     `json:"email"`
	Username      string     `json:"username"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	FullName      string     `json:"full_name"`
	Mobile        string     `json:"mobile,omitempty"`
	UserImage     string     `json:"user_image,omitempty"`
	Roles         []string   `json:"roles"`
	LastLogin     *time.Time `json:"last_login,omitempty"`
}

// Login authenticates a user and creates a session
func (s *AuthService) Login(ctx context.Context, req LoginRequest, ipAddress, userAgent string) (*LoginResponse, error) {
	// Find user by email
	var user base.User
	if err := s.db.Preload("Roles").Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check if user is enabled
	if !user.Enabled {
		return nil, ErrUserDisabled
	}

	// Check if user is locked
	if user.IsLocked() {
		return nil, ErrUserLocked
	}

	// Verify password
	if !user.CheckPassword(req.Password) {
		// Increment failed login attempts
		user.IncrementFailedLogins()
		s.db.Save(&user)
		return nil, ErrInvalidCredentials
	}

	// Generate session ID
	sessionID, err := GenerateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	// Get user roles
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Name
	}

	// Generate JWT token pair
	tokenPair, err := s.jwtService.GenerateTokenPair(user.ID, user.Email, user.Username, roles, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Create session record
	session := base.Session{
		UserID:       user.ID,
		SessionID:    sessionID,
		RefreshToken: tokenPair.RefreshToken,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		LastActivity: time.Now(),
		ExpiresAt:    time.Now().Add(s.config.JWT.RefreshTokenTTL),
		Active:       true,
	}

	if err := s.db.Create(&session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Update user's last login
	user.UpdateLastLogin(ipAddress)
	s.db.Save(&user)

	return &LoginResponse{
		User: &UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			FullName:  user.FullName,
			Mobile:    user.Mobile,
			UserImage: user.UserImage,
			Roles:     roles,
			LastLogin: user.LastLogin,
		},
		TokenPair: tokenPair,
		Message:   "Logged In",
		FullName:  user.FullName,
		Home_Page: "/app/home",
	}, nil
}

// Logout logs out a user and revokes the session
func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	var session base.Session
	if err := s.db.Where("session_id = ?", sessionID).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSessionNotFound
		}
		return err
	}

	// Mark session as logged out
	session.Logout()
	return s.db.Save(&session).Error
}

// RefreshToken refreshes an access token using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// Validate refresh token
	claims, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Find session
	var session base.Session
	if err := s.db.Where("session_id = ?", claims.SessionID).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	// Validate session
	if !session.IsValid() {
		if session.IsExpired() {
			return nil, ErrSessionExpired
		}
		if session.IsRevoked {
			return nil, ErrSessionRevoked
		}
		return nil, errors.New("session is not valid")
	}

	// Generate new token pair
	tokenPair, err := s.jwtService.GenerateTokenPair(
		claims.UserID,
		claims.Email,
		claims.Username,
		claims.Roles,
		claims.SessionID,
	)
	if err != nil {
		return nil, err
	}

	// Update session with new refresh token
	session.RefreshToken = tokenPair.RefreshToken
	session.UpdateActivity()
	s.db.Save(&session)

	return tokenPair, nil
}

// ValidateSession validates a session by session ID
func (s *AuthService) ValidateSession(ctx context.Context, sessionID string) (*base.Session, error) {
	var session base.Session
	if err := s.db.Preload("User.Roles").Where("session_id = ?", sessionID).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	// Validate session
	if !session.IsValid() {
		if session.IsExpired() {
			return nil, ErrSessionExpired
		}
		if session.IsRevoked {
			return nil, ErrSessionRevoked
		}
		return nil, errors.New("session is not valid")
	}

	// Update last activity
	session.UpdateActivity()
	s.db.Save(&session)

	return &session, nil
}

// GetUserByID retrieves a user by ID with roles
func (s *AuthService) GetUserByID(ctx context.Context, userID uint) (*base.User, error) {
	var user base.User
	if err := s.db.Preload("Roles").Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// RevokeSession revokes a session
func (s *AuthService) RevokeSession(ctx context.Context, sessionID, reason string) error {
	var session base.Session
	if err := s.db.Where("session_id = ?", sessionID).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSessionNotFound
		}
		return err
	}

	session.Revoke(reason)
	return s.db.Save(&session).Error
}

// RevokeAllUserSessions revokes all sessions for a user
func (s *AuthService) RevokeAllUserSessions(ctx context.Context, userID uint, reason string) error {
	return s.db.Model(&base.Session{}).
		Where("user_id = ? AND active = ?", userID, true).
		Updates(map[string]interface{}{
			"is_revoked":   true,
			"revoked_at":   time.Now(),
			"revoke_reason": reason,
			"active":       false,
		}).Error
}

// CleanupExpiredSessions removes expired sessions
func (s *AuthService) CleanupExpiredSessions(ctx context.Context) error {
	return s.db.Where("expires_at < ?", time.Now()).Delete(&base.Session{}).Error
}
