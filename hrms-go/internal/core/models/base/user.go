package base

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a system user
type User struct {
	gorm.Model
	Email          string    `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Username       string    `gorm:"uniqueIndex;size:150" json:"username"`
	Password       string    `gorm:"not null" json:"-"` // Never return password in JSON
	FirstName      string    `gorm:"size:150" json:"first_name"`
	LastName       string    `gorm:"size:150" json:"last_name"`
	FullName       string    `gorm:"size:300" json:"full_name"`
	Mobile         string    `gorm:"size:20" json:"mobile"`
	Phone          string    `gorm:"size:20" json:"phone"`
	Gender         string    `gorm:"size:20" json:"gender"`
	BirthDate      *time.Time `json:"birth_date"`
	Location       string    `gorm:"size:255" json:"location"`
	Bio            string    `gorm:"type:text" json:"bio"`
	UserImage      string    `gorm:"type:text" json:"user_image"`

	// Status fields
	Enabled        bool      `gorm:"default:true" json:"enabled"`
	EmailVerified  bool      `gorm:"default:false" json:"email_verified"`
	MobileVerified bool      `gorm:"default:false" json:"mobile_verified"`
	LastLogin      *time.Time `json:"last_login"`
	LastIP         string    `gorm:"size:45" json:"last_ip"`

	// Security fields
	FailedLoginAttempts int       `gorm:"default:0" json:"-"`
	LockedUntil        *time.Time `json:"-"`
	PasswordChangedAt  *time.Time `json:"-"`
	ResetPasswordToken string    `gorm:"size:255" json:"-"`
	ResetPasswordExpiry *time.Time `json:"-"`

	// Preferences
	Language       string `gorm:"size:10;default:'en'" json:"language"`
	TimeZone       string `gorm:"size:50;default:'UTC'" json:"time_zone"`
	DateFormat     string `gorm:"size:20" json:"date_format"`
	TimeFormat     string `gorm:"size:20" json:"time_format"`

	// API Key
	APIKey         string `gorm:"size:255;index" json:"-"`
	APISecret      string `gorm:"size:255" json:"-"`

	// Relationships
	Roles          []Role    `gorm:"many2many:user_roles;" json:"roles"`
	Sessions       []Session `gorm:"foreignKey:UserID" json:"-"`
}

// BeforeCreate hook to hash password and set defaults
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Hash password if not already hashed
	if u.Password != "" && !isHashed(u.Password) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}

	// Set full name
	if u.FullName == "" {
		u.FullName = u.FirstName + " " + u.LastName
	}

	// Set username from email if not provided
	if u.Username == "" {
		u.Username = u.Email
	}

	return nil
}

// BeforeUpdate hook
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// Update full name
	u.FullName = u.FirstName + " " + u.LastName

	return nil
}

// CheckPassword verifies if the provided password matches
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// SetPassword hashes and sets a new password
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	now := time.Now()
	u.PasswordChangedAt = &now
	return nil
}

// IsLocked checks if the user account is locked
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// Lock locks the user account for a duration
func (u *User) Lock(duration time.Duration) {
	lockUntil := time.Now().Add(duration)
	u.LockedUntil = &lockUntil
}

// Unlock unlocks the user account
func (u *User) Unlock() {
	u.LockedUntil = nil
	u.FailedLoginAttempts = 0
}

// IncrementFailedLogins increments failed login attempts
func (u *User) IncrementFailedLogins() {
	u.FailedLoginAttempts++

	// Lock account after 5 failed attempts
	if u.FailedLoginAttempts >= 5 {
		u.Lock(30 * time.Minute)
	}
}

// ResetFailedLogins resets failed login attempts
func (u *User) ResetFailedLogins() {
	u.FailedLoginAttempts = 0
}

// UpdateLastLogin updates the last login timestamp and IP
func (u *User) UpdateLastLogin(ip string) {
	now := time.Now()
	u.LastLogin = &now
	u.LastIP = ip
	u.ResetFailedLogins()
}

// HasRole checks if user has a specific role
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

// HasAnyRole checks if user has any of the specified roles
func (u *User) HasAnyRole(roleNames ...string) bool {
	for _, roleName := range roleNames {
		if u.HasRole(roleName) {
			return true
		}
	}
	return false
}

// IsAdmin checks if user is an administrator
func (u *User) IsAdmin() bool {
	return u.HasAnyRole("System Manager", "Administrator")
}

// Helper function to check if password is already hashed
func isHashed(password string) bool {
	// bcrypt hashes start with $2a$, $2b$, or $2y$
	return len(password) == 60 && password[0] == '$'
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}
