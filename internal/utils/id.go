package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

func NewObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}
