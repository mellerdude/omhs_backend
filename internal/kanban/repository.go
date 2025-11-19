package kanban

import (
	"omhs-backend/internal/requests"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KanbanRepository struct {
	req requests.RequestRepository
}

func NewKanbanRepository(req requests.RequestRepository) *KanbanRepository {
	return &KanbanRepository{req: req}
}

// Load the board JSON for a specific user
func (r *KanbanRepository) GetKanban(userId primitive.ObjectID) (map[string]interface{}, error) {
	doc, err := r.req.Get("data", "Kanbans", userId)

	if err != nil {
		logrus.Warnf("KanbanRepo.GetKanban ERROR for %s → %T: %v",
			userId.Hex(), err, err)
		return nil, err
	}

	return doc.Data, nil
}

// Create a new Kanban document for this user
func (r *KanbanRepository) CreateKanban(doc requests.Document) error {
	return r.req.Create("data", "Kanbans", doc)
}

// Update an existing Kanban document
func (r *KanbanRepository) UpdateKanban(userId primitive.ObjectID, data map[string]interface{}) error {
	return r.req.Update("data", "Kanbans", userId, data)
}
