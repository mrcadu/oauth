package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Profile struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Name        string             `json:"name" binding:"required"`
	Permissions []string           `json:"permissions" binding:"required"`
}
