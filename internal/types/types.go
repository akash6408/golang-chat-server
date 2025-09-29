package types

type User struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Email     string `gorm:"uniqueIndex;null" json:"email,omitempty"`
	Phone     string `gorm:"uniqueIndex;null" json:"phone,omitempty"`
	Password  string `json:"-"` // don’t expose in API response
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type Message struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"` // email or phone of target user
	Message   string `json:"message"`
}

// SignupRequest represents the incoming payload to create a new user account.
type SignupRequest struct {
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Password string `json:"password"`
}

// SignupResponse is returned after successful signup with a token.
type SignupResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	Token string `json:"token"`
}

// LoginRequest represents the login payload for email or phone.
type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Password string `json:"password"`
}

// LoginResponse is returned after successful authentication with a token.
type LoginResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
	Token string `json:"token"`
}
