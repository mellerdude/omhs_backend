package auth

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	FindByUsername(username string) (*User, error)
	FindByEmailAndUsername(email, username string) (*User, error)
	Create(user *User) error
	UpdatePassword(id primitive.ObjectID, newHash string) error
	UpdatePasskey(id primitive.ObjectID, passkey string, at time.Time) error
	InvalidatePasskey(id primitive.ObjectID) error
	UpdateLastLogin(id primitive.ObjectID, at time.Time) error
}

type MongoUserRepository struct {
	client *mongo.Client
}

func NewMongoUserRepository(client *mongo.Client) *MongoUserRepository {
	return &MongoUserRepository{client: client}
}

func (r *MongoUserRepository) collection() *mongo.Collection {
	return r.client.Database("users").Collection("authentication")
}

func (r *MongoUserRepository) FindByUsername(username string) (*User, error) {
	var user User
	err := r.collection().FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	return &user, err
}

func (r *MongoUserRepository) FindByEmailAndUsername(email, username string) (*User, error) {
	var user User
	err := r.collection().FindOne(context.TODO(), bson.M{"email": email, "username": username}).Decode(&user)
	return &user, err
}

func (r *MongoUserRepository) Create(user *User) error {
	_, err := r.collection().InsertOne(context.TODO(), user)
	return err
}

func (r *MongoUserRepository) UpdatePassword(id primitive.ObjectID, newHash string) error {
	_, err := r.collection().UpdateOne(context.TODO(), bson.M{"_id": id},
		bson.M{"$set": bson.M{"password": newHash}, "$unset": bson.M{"passkey": "", "passkeyGeneratedAt": ""}})
	return err
}

func (r *MongoUserRepository) UpdatePasskey(id primitive.ObjectID, passkey string, at time.Time) error {
	_, err := r.collection().UpdateOne(context.TODO(), bson.M{"_id": id},
		bson.M{"$set": bson.M{"passkey": passkey, "passkeyGeneratedAt": at}})
	return err
}

func (r *MongoUserRepository) InvalidatePasskey(id primitive.ObjectID) error {
	_, err := r.collection().UpdateOne(context.TODO(), bson.M{"_id": id},
		bson.M{"$set": bson.M{"passkey": "NOT_PASSKEY"}})
	return err
}

func (r *MongoUserRepository) UpdateLastLogin(id primitive.ObjectID, at time.Time) error {
	_, err := r.collection().UpdateOne(
		context.TODO(),
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"lastLogin": at}},
	)
	return err
}
