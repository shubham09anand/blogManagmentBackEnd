package services

import (
	"context"

	connection "github.com/shubham09anand/blogManagement/connection"
	response "github.com/shubham09anand/blogManagement/error"
	"github.com/shubham09anand/blogManagement/helper"
	model "github.com/shubham09anand/blogManagement/model"
	jwtToken "github.com/shubham09anand/blogManagement/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServices struct{}

var conn, err = connection.ConnectDB()

var collectionUsers = conn.Client.Database("blogManagement").Collection("users")

func (s *UserServices) Signup(data *model.UserSignup) (*response.ServerErrRes, *response.ServerRes, error) {
	// Check same username already exists
	filter := bson.M{"userName": data.UserName}
	var existingUser model.UserSignup

	err := collectionUsers.FindOne(context.Background(), filter).Decode(&existingUser)
	if err == nil {
		return nil, &response.ServerRes{
			Status:   200,
			Success:  false,
			Response: "user already exists",
			Error:    err,
		}, nil
	} else if err != mongo.ErrNoDocuments {
		return nil, &response.ServerRes{
			Status:   404,
			Success:  false,
			Response: "Somthing Went Wrong",
			Error:    err,
		}, nil
	}

	// Proceed to insert the new user
	result, err := collectionUsers.InsertOne(context.Background(), data)
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

func (s *UserServices) Login(userName, password string) (*response.ServerErrRes, *response.ServerRes, error) {
	// Find user by username
	filter := bson.M{"userName": userName}
	var user model.UserSignup

	err := collectionUsers.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &response.ServerRes{
				Status:   404,
				Success:  false,
				Response: "User not found",
				Error:    nil,
			}, nil
		}
		return &response.ServerErrRes{
			Status:   500,
			Response: "Server Failed",
		}, nil, err
	}

	// Check password
	if user.Password != password {
		return nil, &response.ServerRes{
			Status:   401,
			Success:  false,
			Response: "Wrong Credentials",
			Error:    nil,
		}, nil
	}

	// Generate encrypted token
	var key = []byte("your-32-byte-long-key-for-aes-6!") // Ensure this key is 32 bytes for AES-256
	userID := user.Id.Hex()

	encrypted, err := jwtToken.Encrypt(key, userID)
	if err != nil {
		return nil, &response.ServerRes{
			Status:   500,
			Success:  false,
			Response: "Token Generation Failed",
			Error:    nil,
		}, err
	}

	// Create response login data with token and user ID
	type loginReturnData struct {
		Token string `json:"token"`
		Id    string `json:"id"`
	}

	responseLoginData := &loginReturnData{
		Token: encrypted,
		Id:    user.Id.Hex(),
	}

	// Login successful
	return nil, &response.ServerRes{
		Status:   200,
		Success:  true,
		Response: responseLoginData,
		Error:    nil,
	}, nil
}

// to update password
func (s *UserServices) Setting(userID string, password string) (*response.ServerErrRes, *response.ServerRes, error) {

	id, _, _ := helper.ConvertStringToObjectID(userID)

	// Find the user by _id
	filter := bson.M{"_id": id}

	updateValue := bson.M{
		"$set": bson.M{
			"password": password,
		},
	}

	result, err := collectionUsers.UpdateOne(context.Background(), filter, updateValue)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &response.ServerRes{
				Status:   404,
				Success:  false,
				Response: "User not found",
				Error:    nil,
			}, nil
		}
	}

	// Login successful
	return nil, &response.ServerRes{
		Status:   200,
		Success:  true,
		Response: result,
		Error:    nil,
	}, nil
}

// this is to get all writer for displaying their profile
func (s *UserServices) FetchAllUser(ctx context.Context) (*response.ServerErrRes, *response.ServerRes, error) {
	if err != nil {
		return &response.ServerErrRes{
			Status:   400,
			Response: "Server Failed",
		}, nil, err
	}

	lookupUsersStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "profile"}, {Key: "localField", Value: "_id"}, {Key: "foreignField", Value: "userId"}, {Key: "as", Value: "author"}}}}

	// unwindUsersStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$author"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}}

	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "firstName", Value: 1},
		{Key: "lastName", Value: 1},
		{Key: "userName", Value: 1},
		{Key: "photo", Value: "$author.photo"},
	}}}

	cursor, err := collectionUsers.Aggregate(ctx, mongo.Pipeline{lookupUsersStage, projectStage})
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	// Define a variable to store the results
	var tags []bson.M
	if err := cursor.All(ctx, &tags); err != nil {
		return nil, nil, err
	}

	return nil, &response.ServerRes{
		Status:   200,
		Success:  true,
		Response: tags,
		Error:    nil,
	}, nil
}

// this is for getting all writer name only, for search
func (s *UserServices) GetAllWriterName(ctx context.Context) (*response.ServerErrRes, *response.ServerRes, error) {

	if err != nil {
		return &response.ServerErrRes{
			Status:   400,
			Response: "Server Failed",
		}, nil, err
	}

	lookupUsersStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "profile"}, {Key: "localField", Value: "_id"}, {Key: "foreignField", Value: "userId"}, {Key: "as", Value: "author"}}}}

	unwindUsersStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$author"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}}

	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "firstName", Value: 1},
		{Key: "lastName", Value: 1},
		{Key: "photo", Value: "$author.photo"},
	}}}

	cursor, err := collectionUsers.Aggregate(ctx, mongo.Pipeline{lookupUsersStage, unwindUsersStage, projectStage})
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	// Define a variable to store the results
	var tags []bson.M
	if err := cursor.All(ctx, &tags); err != nil {
		return nil, nil, err
	}

	return nil, &response.ServerRes{
		Status:   200,
		Success:  true,
		Response: tags,
		Error:    nil,
	}, nil
}

// this is for getting requested(only one) writer profile details
func (s *UserServices) GetWriterProfile(ctx context.Context, writerId string) (*response.ServerErrRes, *response.ServerRes, error) {

	if err != nil {
		return &response.ServerErrRes{
			Status:   400,
			Response: "Server Failed",
		}, nil, err
	}

	id, _, err := helper.ConvertStringToObjectID(writerId)

	if err != nil {
		return nil, &response.StringToObjevctIdError, nil
	}

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: id}}}}

	lookupUsersStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "profile"}, {Key: "localField", Value: "_id"}, {Key: "foreignField", Value: "userId"}, {Key: "as", Value: "author"}}}}

	// unwindUsersStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$author"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}}

	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "firstName", Value: 1},
		{Key: "lastName", Value: 1},
		{Key: "userName", Value: 1},
		{Key: "createdAt", Value: 1},
		{Key: "interestedTopics", Value: "$author.interestedTopics"},
		{Key: "pronouns", Value: "$author.pronouns"},
		{Key: "aboutYou", Value: "$author.aboutYou"},
		{Key: "phone", Value: "$author.phone"},
		{Key: "photo", Value: "$author.photo"},
		{Key: "email", Value: "$author.email"},
	}}}

	cursor, err := collectionUsers.Aggregate(ctx, mongo.Pipeline{matchStage, lookupUsersStage, projectStage})
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	// Define a variable to store the results
	var writerProfile []bson.M
	if err := cursor.All(ctx, &writerProfile); err != nil {
		return nil, nil, err
	}

	return nil, &response.ServerRes{
		Status:   200,
		Success:  true,
		Response: writerProfile,
		Error:    nil,
	}, nil
}
