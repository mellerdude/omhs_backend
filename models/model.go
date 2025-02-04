package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Document represents a generic MongoDB document
type Document struct {
	ID   primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	Data map[string]interface{} `json:"data" bson:"data"`
}
