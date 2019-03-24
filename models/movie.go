package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Movie is representation of movie entity in mongodb
type Movie struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name string             `bson:"name,omitempty" json:"name,omitempty"`
	//CoverImage  string             `bson:"cover_image,omitempty" json:"cover_image,omitempty"`
	Description string `bson:"description,omitempty" json:"description,omitempty"`
}
