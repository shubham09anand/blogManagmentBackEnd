package services

import (
	"context"
	"errors"

	"log"

	connection "github.com/shubham09anand/blogManagement/Connection"
	response "github.com/shubham09anand/blogManagement/Error"
	helper "github.com/shubham09anand/blogManagement/Helper"
	model "github.com/shubham09anand/blogManagement/Model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentServices struct{}

var conn_5 *connection.Connection
var collectionComment *mongo.Collection

func init() {
	var err error
	conn_5, err = connection.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	collectionComment = conn_5.Client.Database("blogManagement").Collection("comment")
}

func (s *CommentServices) MakeComment(data *model.Comments) (*response.ServerErrRes, *response.ServerRes, error) {
	if conn_5 == nil {
		return &response.ServerErrRes{
			Status:   400,
			Response: "Server Failed",
		}, nil, errors.New("database connection is not initialized")
	}

	result, err := collectionComment.InsertOne(context.Background(), data)
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

func (s *CommentServices) DeleteComment(commentId string) (*response.ServerErrRes, *response.ServerRes, error) {
	if conn_5 == nil {
		return &response.ServerErrRes{
			Status:   400,
			Response: "Server Failed",
		}, nil, errors.New("database connection is not initialized")
	}

	id, _, _ := helper.ConvertStringToObjectID(commentId)
	filter := bson.M{"_id": id}

	result, err := collectionComment.DeleteOne(context.Background(), filter)
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

func (s *CommentServices) GetComments(ctx context.Context, blogId string) (*response.ServerErrRes, *response.ServerRes, error) {
	if conn_5 == nil {
		return &response.ServerErrRes{
			Status:   400,
			Response: "Server Failed",
		}, nil, errors.New("database connection is not initialized")
	}

	id, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, &response.ServerRes{
			Status:   400,
			Success:  false,
			Response: "Failed to convert string to ObjectID",
			Error:    err,
		}, nil
	}

	matchStageBlogId := bson.D{{Key: "$match", Value: bson.D{{Key: "blogId", Value: id}}}}
	lookupUsersStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "users"}, {Key: "localField", Value: "authorId"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "user"}}}}
	unwindUsersStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$user"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}}
	lookupProfileStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "profile"}, {Key: "localField", Value: "authorId"}, {Key: "foreignField", Value: "userId"}, {Key: "as", Value: "author"}}}}
	unwindProfileStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$author"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 1},
		{Key: "authorId", Value: 1},
		{Key: "blogId", Value: 1},
		{Key: "comment", Value: 1},
		{Key: "createdAt", Value: 1},
		{Key: "firstName", Value: "$user.firstName"},
		{Key: "lastName", Value: "$user.lastName"},
		{Key: "photo", Value: "$author.photo"},
	}}}

	cursor, err := collectionComment.Aggregate(ctx, mongo.Pipeline{matchStageBlogId, lookupUsersStage, unwindUsersStage, lookupProfileStage, unwindProfileStage, projectStage})
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(context.Background(), &results); err != nil {
		return nil, nil, err
	}

	if len(results) == 0 {
		return nil, &response.ServerRes{
			Status:   404,
			Success:  false,
			Response: "No comments found for the given blog ID",
			Error:    nil,
		}, nil
	}

	return nil, &response.ServerRes{
		Status:   200,
		Success:  true,
		Response: results,
		Error:    nil,
	}, nil
}
