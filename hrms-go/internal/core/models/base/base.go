package base

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel contains common fields for all models
// Replicates Frappe's document metadata fields
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Frappe-compatible fields
	Name       string `gorm:"uniqueIndex;size:255" json:"name"` // Frappe uses 'name' as unique identifier
	Owner      string `gorm:"size:255" json:"owner"`            // User who created the document
	ModifiedBy string `gorm:"size:255" json:"modified_by"`      // User who last modified
	DocStatus  int    `gorm:"default:0;index" json:"docstatus"` // 0=Draft, 1=Submitted, 2=Cancelled
	Idx        int    `gorm:"default:0" json:"idx"`             // Index/sort order
}

// BeforeCreate hook to set owner
func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	// Set creation timestamp
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now()
	}

	// Set owner from context if available
	if userID, ok := tx.Statement.Context.Value("user_id").(string); ok && b.Owner == "" {
		b.Owner = userID
	}

	return nil
}

// BeforeUpdate hook to set modified_by
func (b *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	// Set modification timestamp
	b.UpdatedAt = time.Now()

	// Set modified_by from context if available
	if userID, ok := tx.Statement.Context.Value("user_id").(string); ok {
		b.ModifiedBy = userID
	}

	return nil
}

// IsSubmitted checks if document is submitted
func (b *BaseModel) IsSubmitted() bool {
	return b.DocStatus == 1
}

// IsCancelled checks if document is cancelled
func (b *BaseModel) IsCancelled() bool {
	return b.DocStatus == 2
}

// IsDraft checks if document is in draft state
func (b *BaseModel) IsDraft() bool {
	return b.DocStatus == 0
}

// Submit marks document as submitted
func (b *BaseModel) Submit() {
	b.DocStatus = 1
}

// Cancel marks document as cancelled
func (b *BaseModel) Cancel() {
	b.DocStatus = 2
}
