package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID         primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Username   string             `json:"username"`
	Password   string             `json:"password,omitempty" bson:"password,omitempty"`
	ProfileIds []string           `json:"profileIds" bson:"profile_ids"`
}

type UserChangePassword struct {
	User
	NewPassword string `json:"newPassword,omitempty"  bson:"-"`
}
