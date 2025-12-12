package middleware

import (
	"strings"

	"github.com/OtchereDev/hrms-go/internal/core/services/auth"
	"github.com/OtchereDev/hrms-go/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware handles authentication
type AuthMiddleware struct {
	jwtService  *auth.JWTService
	authService *auth.AuthService
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(jwtService *auth.JWTService, authService *auth.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService:  jwtService,
		authService: authService,
	}
}

// Authenticate validates JWT token and sets user in context
func (m *AuthMiddleware) Authenticate(c *fiber.Ctx) error {
	// Get token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return response.Unauthorized(c, "missing authorization header")
	}

	// Extract Bearer token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return response.Unauthorized(c, "invalid authorization header format")
	}

	tokenString := parts[1]

	// Validate token
	claims, err := m.jwtService.ValidateToken(tokenString)
	if err != nil {
		if err == auth.ErrExpiredToken {
			return response.Unauthorized(c, "token has expired")
		}
		return response.Unauthorized(c, "invalid token")
	}

	// Ensure it's an access token
	if claims.TokenType != "access" {
		return response.Unauthorized(c, "invalid token type")
	}

	// Validate session
	session, err := m.authService.ValidateSession(c.Context(), claims.SessionID)
	if err != nil {
		if err == auth.ErrSessionExpired {
			return response.Unauthorized(c, "session has expired")
		}
		if err == auth.ErrSessionRevoked {
			return response.Unauthorized(c, "session has been revoked")
		}
		return response.Unauthorized(c, "invalid session")
	}

	// Set user data in context
	c.Locals("user_id", claims.UserID)
	c.Locals("email", claims.Email)
	c.Locals("username", claims.Username)
	c.Locals("roles", claims.Roles)
	c.Locals("session_id", claims.SessionID)
	c.Locals("user", session.User) // Full user object

	return c.Next()
}

// Optional authentication - doesn't fail if no token is provided
func (m *AuthMiddleware) OptionalAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Next()
	}

	// Try to authenticate but don't fail
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		tokenString := parts[1]
		claims, err := m.jwtService.ValidateToken(tokenString)
		if err == nil && claims.TokenType == "access" {
			session, err := m.authService.ValidateSession(c.Context(), claims.SessionID)
			if err == nil {
				c.Locals("user_id", claims.UserID)
				c.Locals("email", claims.Email)
				c.Locals("username", claims.Username)
				c.Locals("roles", claims.Roles)
				c.Locals("session_id", claims.SessionID)
				c.Locals("user", session.User)
			}
		}
	}

	return c.Next()
}

// GetUserID retrieves user ID from context
func GetUserID(c *fiber.Ctx) uint {
	if userID, ok := c.Locals("user_id").(uint); ok {
		return userID
	}
	return 0
}

// GetEmail retrieves email from context
func GetEmail(c *fiber.Ctx) string {
	if email, ok := c.Locals("email").(string); ok {
		return email
	}
	return ""
}

// GetUsername retrieves username from context
func GetUsername(c *fiber.Ctx) string {
	if username, ok := c.Locals("username").(string); ok {
		return username
	}
	return ""
}

// GetRoles retrieves roles from context
func GetRoles(c *fiber.Ctx) []string {
	if roles, ok := c.Locals("roles").([]string); ok {
		return roles
	}
	return []string{}
}

// GetSessionID retrieves session ID from context
func GetSessionID(c *fiber.Ctx) string {
	if sessionID, ok := c.Locals("session_id").(string); ok {
		return sessionID
	}
	return ""
}

// IsAuthenticated checks if user is authenticated
func IsAuthenticated(c *fiber.Ctx) bool {
	return GetUserID(c) > 0
}
