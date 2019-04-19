package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Movie is representation of movie entity in mongodb
type Movie struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	CreatedDate *time.Time         `bson:"created_date,omitempty" json:"created_date,omitempty"`
	CreatedBy   string             `bson:"created_by,omitempty" json:"created_by,omitempty"`
	UpdatedDate *time.Time         `bson:"updated_date,omitempty" json:"updated_date,omitempty"`
	UpdatedBy   string             `bson:"updated_by,omitempty" json:"updated_by,omitempty"`
}
