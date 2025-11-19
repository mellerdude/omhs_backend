package kanban

import (
	"omhs-backend/internal/requests"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KanbanService struct {
	repo KanbanRepository
}

func NewKanbanService(repo KanbanRepository) *KanbanService {
	return &KanbanService{repo: repo}
}

func (s *KanbanService) GetKanban(userId primitive.ObjectID) (map[string]interface{}, error) {
	data, err := s.repo.GetKanban(userId)

	if err != nil {
		// If the kanban doc does not exist, auto-create default
		if err.Error() == "mongo: no documents in result" {
			defaultData := DefaultKanban()

			_, createErr := s.CreateKanban(userId, defaultData)
			if createErr != nil {
				return nil, createErr
			}

			return defaultData, nil
		}

		// Other real errors
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
