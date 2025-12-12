package middleware

import (
	"github.com/OtchereDev/hrms-go/internal/core/models/base"
	"github.com/OtchereDev/hrms-go/pkg/response"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// PermissionMiddleware handles role-based access control
type PermissionMiddleware struct {
	db *gorm.DB
}

// NewPermissionMiddleware creates a new permission middleware
func NewPermissionMiddleware(db *gorm.DB) *PermissionMiddleware {
	return &PermissionMiddleware{
		db: db,
	}
}

// RequireRole checks if user has any of the specified roles
func (m *PermissionMiddleware) RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoles := GetRoles(c)
		if len(userRoles) == 0 {
			return response.Forbidden(c, "no roles assigned")
		}

		// Check if user has any of the required roles
		for _, requiredRole := range roles {
			for _, userRole := range userRoles {
				if userRole == requiredRole {
					return c.Next()
				}
			}
		}

		return response.Forbidden(c, "insufficient permissions")
	}
}

// RequireAllRoles checks if user has all of the specified roles
func (m *PermissionMiddleware) RequireAllRoles(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoles := GetRoles(c)
		if len(userRoles) == 0 {
			return response.Forbidden(c, "no roles assigned")
		}

		// Check if user has all required roles
		for _, requiredRole := range roles {
			hasRole := false
			for _, userRole := range userRoles {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
			if !hasRole {
				return response.Forbidden(c, "insufficient permissions")
			}
		}

		return c.Next()
	}
}

// RequirePermission checks if user has a specific permission on a doctype
func (m *PermissionMiddleware) RequirePermission(docType string, permType string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoles := GetRoles(c)
		if len(userRoles) == 0 {
			return response.Forbidden(c, "no roles assigned")
		}

		// Check if user has the required permission
		var permissions []base.Permission
		err := m.db.Where("doc_type = ? AND role IN ?", docType, userRoles).Find(&permissions).Error
		if err != nil {
			return response.InternalServerError(c, "failed to check permissions")
		}

		// Check if any of the user's roles has the required permission
		for _, perm := range permissions {
			switch permType {
			case "read":
				if perm.Read {
					return c.Next()
				}
			case "write":
				if perm.Write {
					return c.Next()
				}
			case "create":
				if perm.Create {
					return c.Next()
				}
			case "delete":
				if perm.Delete {
					return c.Next()
				}
			case "submit":
				if perm.Submit {
					return c.Next()
				}
			case "cancel":
				if perm.Cancel {
					return c.Next()
				}
			case "amend":
				if perm.Amend {
					return c.Next()
				}
			}
		}

		return response.Forbidden(c, "insufficient permissions for this operation")
	}
}

// IsAdmin checks if user is an administrator
func IsAdmin(c *fiber.Ctx) bool {
	roles := GetRoles(c)
	for _, role := range roles {
		if role == "System Manager" || role == "Administrator" {
			return true
		}
	}
	return false
}

// IsHRManager checks if user is an HR Manager
func IsHRManager(c *fiber.Ctx) bool {
	roles := GetRoles(c)
	for _, role := range roles {
		if role == "HR Manager" || role == "System Manager" {
			return true
		}
	}
	return false
}

// HasRole checks if user has a specific role
func HasRole(c *fiber.Ctx, roleName string) bool {
	roles := GetRoles(c)
	for _, role := range roles {
		if role == roleName {
			return true
		}
	}
	return false
}

// CheckOwnership checks if the current user owns the document
func (m *PermissionMiddleware) CheckOwnership(ownerField string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get owner from request (could be from params, body, or query)
		owner := c.Params(ownerField)
		if owner == "" {
			owner = c.Query(ownerField)
		}

		userEmail := GetEmail(c)

		// Allow if user is admin or HR manager
		if IsAdmin(c) || IsHRManager(c) {
			return c.Next()
		}

		// Check if user is the owner
		if owner != "" && owner != userEmail {
			return response.Forbidden(c, "you can only access your own records")
		}

		return c.Next()
	}
}
