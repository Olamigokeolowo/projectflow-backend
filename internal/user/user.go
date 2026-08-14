package user

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // never serialized into JSON responses
	CreatedAt    time.Time `json:"created_at"`
}