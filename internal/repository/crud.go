package repository

import "go.mongodb.org/mongo-driver/bson/primitive"

type Repository[T any] interface {
	Create(T) (T, error)
	Get(id primitive.ObjectID) (T, error)
	Update(T) (T, error)
	Delete(id primitive.ObjectID) (primitive.ObjectID, error)
}
