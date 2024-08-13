package services

import (
	"context"

	connection "github.com/shubham09anand/blogManagement/connection"
	response "github.com/shubham09anand/blogManagement/error"
	helper "github.com/shubham09anand/blogManagement/helper"
	model "github.com/shubham09anand/blogManagement/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserProfileServices struct{}

var conn_1, err_1 = connection.ConnectDB()

var collectionProfile = conn_1.Client.Database("blogManagement").Collection("profile")

func (s *UserProfileServices) Profile(data *model.UserProfile) (*response.ServerErrRes, *response.ServerRes, error) {

	if err_1 != nil {
		return &response.ServerErrRes{
			Status:   400,
			Response: "Sever Falied",
		}, nil, err_1
	}

	result, err := collectionProfile.InsertOne(context.Background(), data)

	if err != nil {
		return nil, &response.ServerRes{
			Status:   400,
			Success:  false,
			Response: nil,
			Error:    err,
		}, err
	}

	return nil, &response.ServerRes{
		Status:   200,
		Success:  true,
		Response: result,
		Error:    nil,
	}, nil
}

func (s *UserProfileServices) UpdateProfile(data *model.UserProfile) (*response.ServerErrRes, *response.ServerRes, error) {
	if err_1 != nil {
		return &response.ServerErrRes{
			Status:   400,
			Response: "Server Failed",
		}, nil, err_1
	}

	// Construct filter to find the document by its ID
	filter := bson.M{"_id": data.Id}

	// Construct update operation
	update := bson.M{
		"$set": bson.M{
			"email":            data.Email,
			"phone":            data.Phone,
			"pronouns":         data.Pronouns,
			"interestedTopics": data.InterestedTopics,
			"aboutYou":         data.AboutYou,
		},
	}

	// Perform the update operation
	result, err := collectionProfile.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return nil, &response.ServerRes{
			Status:   400,
			Success:  false,
			Response: nil,
			Error:    err,
		}, err
	}

	return nil, &response.ServerRes{
		Status:   200,
		Success:  true,
		Response: result,
		Error:    nil,
	}, nil
}

func (s *UserProfileServices) GetUserPhoto(userId string) (*response.ServerErrRes, *response.ServerRes, error) {

	if err_1 != nil {
		return &response.ServerErrRes{
			Status:   400,
			Response: "Sever Falied",
		}, nil, err_1
	}

	id, _, err := helper.ConvertStringToObjectID(userId)
	if err != nil {
		return &response.ServerErrRes{
			Status:   400,
			Response: "Invalid user ID",
		}, nil, err
	}

	// Construct filter to find the document by its ID
	filter := bson.M{"userId": id}

	// Define a variable to hold the result
	var result model.UserProfile

	err = collectionProfile.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &response.ServerRes{
				Status:   200,
				Success:  false,
				Response: "user haven't made profile yet",
				Error:    nil,
			}, nil
		}
		return nil, &response.ServerRes{
			Status:   200,
			Success:  false,
			Response: "Falied To get Photo",
			Error:    nil,
		}, nil
	}

	return nil, &response.ServerRes{
		Status:   200,
		Success:  true,
		Response: result.Photo,
		Error:    nil,
	}, nil
}
