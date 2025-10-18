package requests

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RequestService struct {
	repo RequestRepository
}

func NewRequestService(repo RequestRepository) *RequestService {
	return &RequestService{repo: repo}
}

func (s *RequestService) Create(database, collection string, data map[string]interface{}) (*Document, error) {
	if database == "" || collection == "" {
		return nil, errors.New("database and collection are required")
	}

	if inner, ok := data["data"].(map[string]interface{}); ok {
		data = inner
	}

	doc := Document{
		ID:   primitive.NewObjectID(),
		Data: data,
	}

	if err := s.repo.Create(database, collection, doc); err != nil {
		return nil, err
	}

	return &doc, nil
}

func (s *RequestService) Get(database, collection, id string) (*Document, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}
	return s.repo.Get(database, collection, objID)
}

func (s *RequestService) Update(database, collection, id string, data map[string]interface{}) (*Document, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}
	if err := s.repo.Update(database, collection, objID, data); err != nil {
		return nil, err
	}
	return s.repo.Get(database, collection, objID)
}

func (s *RequestService) Delete(database, collection, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id format")
	}
	return s.repo.Delete(database, collection, objID)
}

func (s *RequestService) GetAll(database, collection string) ([]Document, error) {
	return s.repo.GetAll(database, collection)
}
