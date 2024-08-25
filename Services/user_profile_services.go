package services

import (
	"context"
	"fmt"

	connection "github.com/shubham09anand/blogManagement/Connection"
	response "github.com/shubham09anand/blogManagement/Error"
	helper "github.com/shubham09anand/blogManagement/Helper"
	model "github.com/shubham09anand/blogManagement/Model"
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

func (s *UserProfileServices) UpdateProfile(firstName string, lastName string, data *model.UserProfile) (*response.ServerErrRes, *response.ServerRes, error) {
	if err_1 != nil {
		return &response.ServerErrRes{
			Status:   400,
			Response: "Server Failed",
		}, nil, err_1
	}

	session, err := conn_1.Client.StartSession()
	if err != nil {
		return nil, &response.ServerRes{
			Status:   500,
			Success:  false,
			Response: "Failed to start session",
			Error:    err,
		}, err
	}

	fmt.Println("Services")
	fmt.Println(firstName)
	fmt.Println(lastName)
	fmt.Println(data)

	defer session.EndSession(context.Background())

	err = mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		// Update profile document
		filterProfile := bson.M{"userId": data.UserId}
		updateProfile := bson.M{
			"$set": bson.M{
				"email":            data.Email,
				"photo":            data.Photo,
				"phone":            data.Phone,
				"pronouns":         data.Pronouns,
				"interestedTopics": data.InterestedTopics,
				"aboutYou":         data.AboutYou,
			},
		}

		_, err := collectionProfile.UpdateOne(sc, filterProfile, updateProfile)
		if err != nil {
			return err
		}

		// Update user document
		collectionUsers := conn_1.Client.Database("blogManagement").Collection("users")
		filterUser := bson.M{"_id": data.UserId}
		updateUsers := bson.M{
			"$set": bson.M{
				"firstName": firstName,
				"lastName":  lastName,
			},
		}

		_, err = collectionUsers.UpdateOne(sc, filterUser, updateUsers)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return &response.ServerErrRes{
			Status:   500,
			Response: "Failed to update documents",
		}, nil, err
	}

	return nil, &response.ServerRes{
		Status:   200,
		Success:  true,
		Response: "Documents updated successfully",
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

// func (s *UserProfileServices) GetWriterProfile(ctx context.Context, userId string) (*response.ServerErrRes, *response.ServerRes, error) {

// 	if err != nil {
// 		return &response.ServerErrRes{
// 			Status:   400,
// 			Response: "Server Failed",
// 		}, nil, err
// 	}

// 	id, _, err := helper.ConvertStringToObjectID(userId)

// 	if err != nil {
// 		return nil, &response.StringToObjevctIdError, nil
// 	}

// 	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: id}}}}

// 	lookupUsersStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "users"}, {Key: "localField", Value: "userId"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "user"}}}}

// 	// unwindUsersStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$user"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}}

// 	projectStage := bson.D{{Key: "$project", Value: bson.D{
// 		{Key: "phone", Value: 1},
// 		{Key: "photo", Value: 1},
// 		{Key: "email", Value: 1},
// 		{Key: "pronouns", Value: 1},
// 		{Key: "aboutYou", Value: 1},
// 		{Key: "firstName", Value: "$user.firstName"},
// 		{Key: "lastName", Value: "$user.lastName"},
// 		{Key: "email", Value: "$user.email"},
// 	}}}

// 	cursor, err := collectionUsers.Aggregate(ctx, mongo.Pipeline{matchStage, lookupUsersStage, projectStage})
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	// Define a variable to store the results
// 	var tags []bson.M
// 	if err := cursor.All(ctx, &tags); err != nil {
// 		return nil, nil, err
// 	}

// 	return nil, &response.ServerRes{
// 		Status:   200,
// 		Success:  true,
// 		Response: tags,
// 		Error:    nil,
// 	}, nil
// }
