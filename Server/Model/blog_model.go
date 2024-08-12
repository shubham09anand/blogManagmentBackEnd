package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	AuthorId  primitive.ObjectID `json:"authorId" bson:"authorId" binding:"required"`
	Title     string             `json:"title" bson:"title" binding:"required"`
	Tags      []string           `json:"tags" bson:"tags" binding:"required"`
	Content   string             `json:"content" bson:"content" binding:"required"`
	BlogPhoto string             `json:"blogPhoto" bson:"blogPhoto"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt" binding:"required"`
}

type Comments struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	AuthorId  primitive.ObjectID `json:"authorId" bson:"authorId" binding:"required"`
	BlogId    primitive.ObjectID `json:"blogId" bson:"blogId" binding:"required"`
	Comment   string             `json:"comment" bson:"comment" binding:"required"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt" binding:"required"`
}
