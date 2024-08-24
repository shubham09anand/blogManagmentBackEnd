package services

import (
	"context"
	"errors"

	connection "github.com/shubham09anand/blogManagement/connection"
	response "github.com/shubham09anand/blogManagement/error"
	helper "github.com/shubham09anand/blogManagement/helper"
	model "github.com/shubham09anand/blogManagement/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogServices struct{}

var conn_2, err_2 = connection.ConnectDB()

var collectionBlog = conn_2.Client.Database("blogManagement").Collection("blog")

// pipelining done

func (s *BlogServices) CreateBlog(data *model.Blog) (*response.ServerErrRes, *response.ServerRes, error) {
	if err_2 != nil {
		return &response.ServerErrRes{
			Status:   400,
			Response: "Sever Falied",
		}, nil, err_2
	}

	result, err := collectionBlog.InsertOne(context.Background(), data)

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

func (s *BlogServices) UpdateBlog(data *model.Blog) (*response.ServerErrRes, *response.ServerRes, error) {
	if err_2 != nil {
		return &response.ServerErrRes{
			Status:   400,
			Response: "Sever Falied",
		}, nil, err_2
	}

	filter := bson.M{"_id": data.Id, "authorId": data.AuthorId}

	updateValue := bson.M{
		"$set": bson.M{
			"title":     data.Title,
			"content":   data.Content,
			"tags":      data.Tags,
			"blogPhoto": data.BlogPhoto,
		},
	}

	result, err := collectionBlog.UpdateOne(context.Background(), filter, updateValue)

	if err != nil {
		return nil, &response.ServerRes{
			Status:   400,
			Success:  false,
			Response: nil,
			Error:    err,
		}, err
	}

	nodocFound := errors.New("no doument found")

	if result.MatchedCount == 0 {
		return nil, &response.ServerRes{
			Status:   404,
			Success:  false,
			Response: result,
			Error:    err,
		}, nodocFound

	}

	return nil, &response.ServerRes{
		Status:   200,
		Success:  true,
		Response: result,
		Error:    nil,
	}, nil
}

func (s *BlogServices) DeleteBlog(blogId string, blogAuthorId string) (*response.ServerErrRes, *response.ServerRes, error) {
	if err_2 != nil {
		return &response.ServerErrRes{
			Status:   400,
			Response: "Sever Falied",
		}, nil, err_2
	}

	// fmt.Println("This is Delete")
	id, _, _ := helper.ConvertStringToObjectID(blogId)
	authorId, _, _ := helper.ConvertStringToObjectID(blogAuthorId)

	filter := bson.M{"_id": id, "authorId": authorId}

	result, err := collectionBlog.DeleteOne(context.Background(), filter)

	if err != nil {
		return nil, &response.ServerRes{
			Status:   400,
			Success:  false,
			Response: nil,
			Error:    err,
		}, err
	}

	if result.DeletedCount == 0 {
		return &response.ServerErrRes{
			Status:   404,
			Response: "No document found",
		}, nil, errors.New("no document found")
	}

	return nil, &response.ServerRes{
		Status:   200,
		Success:  true,
		Response: result,
		Error:    nil,
	}, nil
}

func (s *BlogServices) FetchAllBlog(ctx context.Context) (*response.ServerErrRes, *response.ServerRes, error) {
	lookupUsersStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "users"}, {Key: "localField", Value: "authorId"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "author"}}}}

	unwindUsersStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$author"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}}

	lookupProfileStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "profile"}, {Key: "localField", Value: "authorId"}, {Key: "foreignField", Value: "userId"}, {Key: "as", Value: "profile"}}}}

	unwindProfileStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$profile"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}}

	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 1},
		{Key: "title", Value: 1},
		{Key: "tags", Value: 1},
		{Key: "content", Value: 1}, // Fetch the entire content field
		{Key: "blogPhoto", Value: 1},
		{Key: "createdAt", Value: 1},
		{Key: "authorId", Value: 1},
		{Key: "firstName", Value: "$author.firstName"},
		{Key: "lastName", Value: "$author.lastName"},
		{Key: "photo", Value: "$profile.photo"},
	}}}

	// Execute the aggregation pipeline
	cursor, err := collectionBlog.Aggregate(ctx, mongo.Pipeline{lookupUsersStage, unwindUsersStage, lookupProfileStage, unwindProfileStage, projectStage})
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, nil, err
	}

	// Process the results to strip HTML tags and truncate the content field
	for i, result := range results {
		if content, ok := result["content"].(string); ok {
			plainText, err := helper.ExtractPlainText(content)
			if err != nil {
				return nil, nil, err
			}
			if len(plainText) > 250 {
				plainText = plainText[:250] // Truncate to 250 characters
			}
			results[i]["content"] = plainText
		}
	}

	return nil, &response.ServerRes{
		Status:   200,
		Success:  true,
		Response: results,
		Error:    nil,
	}, nil
}

func (s *BlogServices) FetchOneBlog(ctx context.Context, blogId string) (*response.ServerErrRes, *response.ServerRes, error) {

	if err_2 != nil {
		return &response.ServerErrRes{
			Status:   400,
			Response: "Sever Falied",
		}, nil, err_2
	}

	id, _, _ := helper.ConvertStringToObjectID(blogId)

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: id}}}}

	lookupUsersStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "users"}, {Key: "localField", Value: "authorId"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "author"}}}}

	unwindUsersStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$author"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}}

	lookupProfileStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "profile"}, {Key: "localField", Value: "authorId"}, {Key: "foreignField", Value: "userId"}, {Key: "as", Value: "profile"}}}}

	unwindProfileStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$profile"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}}

	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 1},
		{Key: "title", Value: 1},
		{Key: "tags", Value: 1},
		{Key: "content", Value: 1},
		{Key: "blogPhoto", Value: 1},
		{Key: "createdAt", Value: 1},
		{Key: "authorId", Value: 1},
		{Key: "firstName", Value: "$author.firstName"},
		{Key: "lastName", Value: "$author.lastName"},
		{Key: "photo", Value: "$profile.photo"},
	}}}

	// Execute the aggregation pipeline
	cursor, err := collectionBlog.Aggregate(ctx, mongo.Pipeline{matchStage, lookupUsersStage, unwindUsersStage, lookupProfileStage, unwindProfileStage, projectStage})
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, nil, err
	}

	// Debug: Print the fetched results to check the data structure
	// fmt.Println("Fetched blogs with authors:", results)

	return nil, &response.ServerRes{
		Status:   200,
		Success:  true,
		Response: results,
		Error:    nil,
	}, nil
}

// FetchWriterBlog retrieves blog posts by a specific writer and processes the content field to remove HTML tags
func (s *BlogServices) FetchWriterBlog(ctx context.Context, writerId string) (*response.ServerErrRes, *response.ServerRes, error) {
	authorId, _, err := helper.ConvertStringToObjectID(writerId)
	if err != nil {
		return &response.ServerErrRes{
			Status:   400,
			Response: "Invalid writer ID",
		}, nil, err
	}

	matchStageAuthorId := bson.D{{Key: "$match", Value: bson.D{{Key: "authorId", Value: authorId}}}}

	lookupUsersStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "users"}, {Key: "localField", Value: "authorId"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "author"}}}}

	unwindUsersStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$author"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}}

	lookupProfileStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "profile"}, {Key: "localField", Value: "authorId"}, {Key: "foreignField", Value: "userId"}, {Key: "as", Value: "profile"}}}}

	unwindProfileStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$profile"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}}

	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 1},
		{Key: "title", Value: 1},
		{Key: "tags", Value: 1},
		{Key: "content", Value: 1}, // Fetch the entire content field
		{Key: "createdAt", Value: 1},
		{Key: "authorId", Value: 1},
		{Key: "blogPhoto", Value: 1},
		{Key: "firstName", Value: "$author.firstName"},
		{Key: "lastName", Value: "$author.lastName"},
		{Key: "photo", Value: "$profile.photo"},
	}}}

	// Execute the aggregation pipeline
	cursor, err := collectionBlog.Aggregate(ctx, mongo.Pipeline{matchStageAuthorId, lookupUsersStage, unwindUsersStage, lookupProfileStage, unwindProfileStage, projectStage})
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, nil, err
	}

	// Process the results to strip HTML tags and truncate the content field
	for i, result := range results {
		if content, ok := result["content"].(string); ok {
			plainText, err := helper.ExtractPlainText(content)
			if err != nil {
				return nil, nil, err
			}
			if len(plainText) > 250 {
				plainText = plainText[:250] // Truncate to 250 characters
			}
			results[i]["content"] = plainText
		}
	}

	if len(results) == 0 {
		return nil, &response.ServerRes{
			Status:   200,
			Success:  false,
			Response: "No blog found with the given writer ID",
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

func (s *BlogServices) SearchBlogByTitile(ctx context.Context, blogTitle string) (*response.ServerErrRes, *response.ServerRes, error) {

	if err_2 != nil {
		return &response.ServerErrRes{
			Status:   400,
			Response: "Sever Falied",
		}, nil, err_2
	}

	matchStageAuthorId := bson.D{{Key: "$match", Value: bson.D{{Key: "title", Value: blogTitle}}}}

	lookupUsersStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "users"}, {Key: "localField", Value: "authorId"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "author"}}}}

	unwindUsersStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$author"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}}

	lookupProfileStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "profile"}, {Key: "localField", Value: "authorId"}, {Key: "foreignField", Value: "userId"}, {Key: "as", Value: "profile"}}}}

	unwindProfileStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$profile"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}}

	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 1},
		{Key: "title", Value: 1},
		{Key: "tags", Value: 1},
		{Key: "content", Value: 1},
		{Key: "createdAt", Value: 1},
		{Key: "authorId", Value: 1},
		{Key: "firstName", Value: "$author.firstName"},
		{Key: "lastName", Value: "$author.lastName"},
		{Key: "photo", Value: "$profile.photo"},
	}}}

	// Execute the aggregation pipeline
	cursor, err := collectionBlog.Aggregate(ctx, mongo.Pipeline{matchStageAuthorId, lookupUsersStage, unwindUsersStage, lookupProfileStage, unwindProfileStage, projectStage})
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, nil, err
	}

	if err == mongo.ErrNoDocuments {
		return nil, &response.ServerRes{
			Status:   404,
			Success:  false,
			Response: "No blog found with the given ID",
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

func (s *BlogServices) SearchBlogByTag(ctx context.Context, tag string) (*response.ServerErrRes, *response.ServerRes, error) {

	if err_2 != nil {
		return &response.ServerErrRes{
			Status:   400,
			Response: "Sever Falied",
		}, nil, err_2
	}

	matchStageAuthorId := bson.D{{Key: "$match", Value: bson.D{{Key: "tags", Value: tag}}}}

	lookupUsersStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "users"}, {Key: "localField", Value: "authorId"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "author"}}}}

	unwindUsersStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$author"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}}

	lookupProfileStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "profile"}, {Key: "localField", Value: "authorId"}, {Key: "foreignField", Value: "userId"}, {Key: "as", Value: "profile"}}}}

	unwindProfileStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$profile"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}}

	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 1},
		{Key: "title", Value: 1},
		{Key: "tags", Value: 1},
		{Key: "content", Value: 1},
		{Key: "createdAt", Value: 1},
		{Key: "authorId", Value: 1},
		{Key: "blogPhoto", Value: 1},
		{Key: "createdAt", Value: 1},
		{Key: "firstName", Value: "$author.firstName"},
		{Key: "lastName", Value: "$author.lastName"},
		{Key: "photo", Value: "$profile.photo"},
	}}}

	// Execute the aggregation pipeline
	cursor, err := collectionBlog.Aggregate(ctx, mongo.Pipeline{matchStageAuthorId, lookupUsersStage, unwindUsersStage, lookupProfileStage, unwindProfileStage, projectStage})
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, nil, err
	}

	if err == mongo.ErrNoDocuments {
		return nil, &response.ServerRes{
			Status:   404,
			Success:  false,
			Response: "No blog found with the given ID",
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

func (s *BlogServices) SearchBlogByWriterName(ctx context.Context, tag string) (*response.ServerErrRes, *response.ServerRes, error) {

	if err_2 != nil {
		return &response.ServerErrRes{
			Status:   400,
			Response: "Sever Falied",
		}, nil, err_2
	}

	matchStageAuthorId := bson.D{{Key: "$match", Value: bson.D{{Key: "tags", Value: tag}}}}

	lookupUsersStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "users"}, {Key: "localField", Value: "authorId"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "author"}}}}

	unwindUsersStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$author"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}}

	lookupProfileStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "profile"}, {Key: "localField", Value: "authorId"}, {Key: "foreignField", Value: "userId"}, {Key: "as", Value: "profile"}}}}

	unwindProfileStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$profile"}, {Key: "preserveNullAndEmptyArrays", Value: false}}}}

	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 1},
		{Key: "title", Value: 1},
		{Key: "tags", Value: 1},
		{Key: "content", Value: 1},
		{Key: "createdAt", Value: 1},
		{Key: "authorId", Value: 1},
		{Key: "firstName", Value: "$author.firstName"},
		{Key: "lastName", Value: "$author.lastName"},
		{Key: "photo", Value: "$profile.photo"},
	}}}

	// Execute the aggregation pipeline
	cursor, err := collectionBlog.Aggregate(ctx, mongo.Pipeline{matchStageAuthorId, lookupUsersStage, unwindUsersStage, lookupProfileStage, unwindProfileStage, projectStage})
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, nil, err
	}

	if err == mongo.ErrNoDocuments {
		return nil, &response.ServerRes{
			Status:   404,
			Success:  false,
			Response: "No blog found with the given ID",
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

func (s *BlogServices) GetAllTags(ctx context.Context) (*response.ServerErrRes, *response.ServerRes, error) {

	if err_2 != nil {
		return &response.ServerErrRes{
			Status:   400,
			Response: "Server Failed",
		}, nil, err_2
	}

	projectStage := bson.D{{Key: "$project", Value: bson.D{{Key: "tags", Value: 1}}}}

	cursor, err := collectionBlog.Aggregate(ctx, mongo.Pipeline{projectStage})
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
