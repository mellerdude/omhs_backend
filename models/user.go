package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system.
type User struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	Username           string             `json:"username"`
	Password           string             `json:"password"`
	Email              string             `json:"email"`
	IsAdmin            bool               `json:"isAdmin"`
	Token              string             `json:"token"`
	LastLogin          time.Time          `json:"lastLogin"`
	Passkey            string             `json:"passkey"`
	PasskeyGeneratedAt time.Time          `json:"passkeyGeneratedAt"`
}
