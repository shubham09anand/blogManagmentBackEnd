package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSignup struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName  string             `json:"userName" bson:"userName" binding:"required"`
	FirstName string             `json:"firstName" bson:"firstName" binding:"required"`
	LastName  string             `json:"lastName" bson:"lastName" binding:"required"`
	Password  string             `json:"password" bson:"password" binding:"required"`
	CreatedAt string             `json:"createdAt" bson:"createdAt" binding:"required"`
}

type UserProfile struct {
	Id               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId           primitive.ObjectID `json:"userId" bson:"userId" binding:"required"`
	Email            string             `json:"email" bson:"email" binding:"required"`
	Photo            string             `json:"photo" bson:"photo"`
	Phone            string             `json:"phone" bson:"phone"`
	Pronouns         string             `json:"pronouns" bson:"pronouns" binding:"required"`
	InterestedTopics string             `json:"interestedTopics" bson:"interestedTopics" binding:"required"`
	AboutYou         string             `json:"aboutYou" bson:"aboutYou"`
}
