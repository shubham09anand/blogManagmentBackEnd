package helper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	response "github.com/shubham09anand/blogManagement/error"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/html"
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

// get text out of raw html
func ExtractPlainText(htmlContent string) (string, error) {
	// Parse the HTML document
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}

	// Function to extract text from nodes
	var extractText func(*html.Node)
	var textContent strings.Builder
	extractText = func(n *html.Node) {
		if n.Type == html.TextNode {
			textContent.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}
	extractText(doc)

	return textContent.String(), nil
}
