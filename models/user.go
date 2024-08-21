package models

import "time"

// User represents a user entity.
type User struct {
	ID              int       `json:"id"`
	Username        string    `json:"username"`
	Password        string    `json:"password"`
	Email           string    `json:"email"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Url             string    `json:"url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CurrentPassword string    `json:"current_password,omitempty"`
	NewPassword     string    `json:"new_password,omitempty"`
}
