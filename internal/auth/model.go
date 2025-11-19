package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system.
type User struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username           string             `bson:"username" json:"username"`
	Password           string             `bson:"password" json:"password"`
	Email              string             `bson:"email" json:"email"`
	IsAdmin            bool               `bson:"isAdmin" json:"isAdmin"`
	LastLogin          time.Time          `bson:"lastLogin" json:"lastLogin"`
	Passkey            string             `bson:"passkey" json:"passkey"`
	PasskeyGeneratedAt time.Time          `bson:"passkeyGeneratedAt" json:"passkeyGeneratedAt"`
}

// DTOs (request payloads)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResetPasswordRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type ChangePasswordRequest struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	Passkey     string `json:"passkey"`
	NewPassword string `json:"newPassword"`
}
