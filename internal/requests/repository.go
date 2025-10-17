package requests

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RequestRepository interface {
	Create(database, collection string, doc Document) error
	Get(database, collection string, id primitive.ObjectID) (*Document, error)
	Update(database, collection string, id primitive.ObjectID, data map[string]interface{}) error
	Delete(database, collection string, id primitive.ObjectID) error
	GetAll(database, collection string) ([]Document, error)
}

type MongoRequestRepository struct {
	client *mongo.Client
}

func NewMongoRequestRepository(client *mongo.Client) *MongoRequestRepository {
	return &MongoRequestRepository{client: client}
}

func (r *MongoRequestRepository) col(database, collection string) *mongo.Collection {
	return r.client.Database(database).Collection(collection)
}

func (r *MongoRequestRepository) Create(database, collection string, doc Document) error {
	_, err := r.col(database, collection).InsertOne(context.TODO(), doc)
	logrus.Infof("Document before insert: %+v", doc)
	return err
}

func (r *MongoRequestRepository) Get(database, collection string, id primitive.ObjectID) (*Document, error) {
	var raw bson.M
	err := r.col(database, collection).FindOne(context.TODO(), bson.M{"_id": id}).Decode(&raw)
	if err != nil {
		return nil, err
	}

	data, ok := raw["data"].(map[string]interface{})
	if !ok {
		// fallback: everything except _id becomes data
		delete(raw, "_id")
		data = raw
	}

	doc := &Document{
		ID:   id,
		Data: data,
	}

	logrus.Infof("Decoded Document: %+v", doc)
	return doc, nil
}

func (r *MongoRequestRepository) Update(database, collection string, id primitive.ObjectID, data map[string]interface{}) error {
	_, err := r.col(database, collection).UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": data})
	return err
}

func (r *MongoRequestRepository) Delete(database, collection string, id primitive.ObjectID) error {
	_, err := r.col(database, collection).DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}

func (r *MongoRequestRepository) GetAll(database, collection string) ([]Document, error) {
	cursor, err := r.col(database, collection).Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var docs []Document
	for cursor.Next(context.TODO()) {
		var d Document
		if err := cursor.Decode(&d); err != nil {
			return nil, err
		}
		docs = append(docs, d)
	}
	return docs, cursor.Err()
}
