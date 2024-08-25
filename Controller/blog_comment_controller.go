package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	helper "github.com/shubham09anand/blogManagement/Helper"
	model "github.com/shubham09anand/blogManagement/Model"
	services "github.com/shubham09anand/blogManagement/Services"
)

type MakeCommentController struct {
	Services *services.CommentServices
}

type DeleteCommentController struct {
	Services *services.CommentServices
}

type GetCommentController struct {
	Services *services.CommentServices
}

func (s *MakeCommentController) MakeComment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var commentData model.Comments

		if !(helper.BindJSON(ctx, &commentData)) {
			return
		}

		_, response, err := s.Services.MakeComment(&commentData)

		if !(helper.InternalServerError(ctx, err)) {
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"response": response})
	}
}

func (s *DeleteCommentController) DeleteComment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var deleteCommentModel struct {
			ID string `json:"_id" binding:"required"`
		}

		deleteCommentData := &deleteCommentModel

		_, response, err := s.Services.DeleteComment(deleteCommentData.ID)
		if !(helper.InternalServerError(ctx, err)) {
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"response": response})
	}
}

func (s *GetCommentController) GetComments() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		Id := ctx.Param("blogId")

		_, response, err := s.Services.GetComments(ctx, Id)
		if !(helper.InternalServerError(ctx, err)) {
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"response": response})
	}
}
