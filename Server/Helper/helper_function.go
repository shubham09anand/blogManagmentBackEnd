package helper

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	response "github.com/shubham09anand/blogManagement/error"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BindJSON(ctx *gin.Context, obj interface{}) bool {
	if err := ctx.ShouldBindJSON(obj); err != nil {
		bindingError := response.BindingErr
		bindingError.Error = fmt.Sprintf("Invalid request format: %s", err.Error())
		ctx.JSON(http.StatusBadRequest, bindingError)
		return false
	}
	return true
}

func InternalServerError(ctx *gin.Context, err error) bool {
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data", "details": err.Error()})
		return false
	}
	return true
}

func ConvertStringToObjectID(inputID string) (primitive.ObjectID, *response.ServerRes, error) {
	fmt.Println("hi")
	id, err := primitive.ObjectIDFromHex(inputID)
	if err != nil {
		return primitive.NilObjectID, &response.ServerRes{
			Status:   400,
			Success:  false,
			Response: "Failed to convert string to ObjectID",
			Error:    err,
		}, err
	}

	return id, nil, nil
}
