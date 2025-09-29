package models

// User represents a user in the system
type User struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Email     string `gorm:"uniqueIndex;null" json:"email,omitempty"`
	Phone     string `gorm:"uniqueIndex;null" json:"phone,omitempty"`
	Password  string `json:"-"` // don’t expose in API response
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
