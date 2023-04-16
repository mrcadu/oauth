package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	Username string             `json:"username"`
	Password string             `json:"password"`
}
