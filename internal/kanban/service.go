package kanban

import (
	"errors"
	"omhs-backend/internal/requests"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type KanbanService struct {
	repo KanbanRepository
}

func NewKanbanService(repo KanbanRepository) *KanbanService {
	return &KanbanService{repo: repo}
}

func (s *KanbanService) GetKanban(userId primitive.ObjectID) (map[string]interface{}, error) {
	data, err := s.repo.GetKanban(userId)
	logrus.Warnf("KanbanService.GetKanban → repo returned error: %T %v", err, err)

	if err != nil {
		// Document doesn't exist → create default
		if errors.Is(err, mongo.ErrNoDocuments) {
			defaultData := DefaultKanban()

			_, createErr := s.CreateKanban(userId, defaultData)
			if createErr != nil {
				return nil, createErr
			}

			return defaultData, nil
		}

		// Other errors
		return nil, err
	}

	return data, nil
}

func (s *KanbanService) CreateKanban(userId primitive.ObjectID, data map[string]interface{}) (*requests.Document, error) {
	doc := requests.Document{
		ID:   userId,
		Data: data,
	}

	if err := s.repo.CreateKanban(doc); err != nil {
		return nil, err
	}

	return &doc, nil
}

func (s *KanbanService) UpdateKanban(userId primitive.ObjectID, data map[string]interface{}) error {
	return s.repo.UpdateKanban(userId, data)
}
