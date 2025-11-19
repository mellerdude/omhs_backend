package kanban

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Kanban struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty"`
	Data      map[string]interface{} `bson:"data"`
	UpdatedAt time.Time              `bson:"updatedAt"`
}
