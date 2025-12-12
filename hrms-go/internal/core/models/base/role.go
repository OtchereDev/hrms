package base

import "gorm.io/gorm"

// Role represents a user role
type Role struct {
	gorm.Model
	Name        string       `gorm:"uniqueIndex;not null;size:150" json:"name"`
	Description string       `gorm:"type:text" json:"description"`
	Disabled    bool         `gorm:"default:false" json:"disabled"`
	IsDeskUser  bool         `gorm:"default:true" json:"is_desk_user"` // Can access desk/admin interface

	// Relationships
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Users       []User       `gorm:"many2many:user_roles;" json:"-"`
}

// HasPermission checks if role has a specific permission
func (r *Role) HasPermission(docType string, permType string) bool {
	for _, perm := range r.Permissions {
		if perm.DocType == docType {
			switch permType {
			case "read":
				return perm.Read
			case "write":
				return perm.Write
			case "create":
				return perm.Create
			case "delete":
				return perm.Delete
			case "submit":
				return perm.Submit
			case "cancel":
				return perm.Cancel
			case "amend":
				return perm.Amend
			}
		}
	}
	return false
}

// TableName specifies the table name for Role model
func (Role) TableName() string {
	return "roles"
}

// Permission represents a document-level permission
type Permission struct {
	gorm.Model
	Role     string `gorm:"index;not null;size:150" json:"role"`
	DocType  string `gorm:"index;not null;size:150" json:"doc_type"` // The model/document type

	// Permission types (Frappe-compatible)
	Read     bool `gorm:"default:false" json:"read"`
	Write    bool `gorm:"default:false" json:"write"`
	Create   bool `gorm:"default:false" json:"create"`
	Delete   bool `gorm:"default:false" json:"delete"`
	Submit   bool `gorm:"default:false" json:"submit"`
	Cancel   bool `gorm:"default:false" json:"cancel"`
	Amend    bool `gorm:"default:false" json:"amend"`
	Print    bool `gorm:"default:false" json:"print"`
	Email    bool `gorm:"default:false" json:"email"`
	Report   bool `gorm:"default:false" json:"report"`
	Import   bool `gorm:"default:false" json:"import"`
	Export   bool `gorm:"default:false" json:"export"`
	Share    bool `gorm:"default:false" json:"share"`

	// Advanced permissions
	IfOwner  bool   `gorm:"default:false" json:"if_owner"`  // Only if user is the owner
	ApplyUserPermissions bool `gorm:"default:false" json:"apply_user_permissions"`
	PermLevel int    `gorm:"default:0" json:"perm_level"`   // Permission level (0, 1, 2...)

	// Conditions
	// MatchField string `gorm:"size:150" json:"match_field"` // Field to match (e.g., "department")
}

// TableName specifies the table name for Permission model
func (Permission) TableName() string {
	return "permissions"
}

// DefaultRoles returns default roles for the system
func DefaultRoles() []Role {
	return []Role{
		{
			Name:        "System Manager",
			Description: "Full system access",
			IsDeskUser:  true,
		},
		{
			Name:        "HR Manager",
			Description: "HR management with full access to HR and Payroll modules",
			IsDeskUser:  true,
		},
		{
			Name:        "HR User",
			Description: "Basic HR operations",
			IsDeskUser:  true,
		},
		{
			Name:        "Employee",
			Description: "Employee self-service access",
			IsDeskUser:  false,
		},
		{
			Name:        "Leave Approver",
			Description: "Can approve/reject leave applications",
			IsDeskUser:  true,
		},
		{
			Name:        "Expense Approver",
			Description: "Can approve/reject expense claims",
			IsDeskUser:  true,
		},
		{
			Name:        "Shift Request Approver",
			Description: "Can approve/reject shift requests",
			IsDeskUser:  true,
		},
		{
			Name:        "Guest",
			Description: "Limited read-only access",
			IsDeskUser:  false,
		},
	}
}
