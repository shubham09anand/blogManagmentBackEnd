package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	helper "github.com/shubham09anand/blogManagement/helper"
	model "github.com/shubham09anand/blogManagement/model"
	services "github.com/shubham09anand/blogManagement/services"
)

type BlogController struct {
	Services *services.BlogServices
}

type DeleteBlogController struct {
	Services *services.BlogServices
}

type GetAllBlogController struct {
	Services *services.BlogServices
}

type GetOneBlogController struct {
	Services *services.BlogServices
}

type GetWriterBlogController struct {
	Services *services.BlogServices
}

type SearchBlogByTitleController struct {
	Services *services.BlogServices
}

type SearchBlogByTagController struct {
	Services *services.BlogServices
}

type GetAllTagController struct {
	Services *services.BlogServices
}

func (s *BlogController) CreateBlog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var blogData model.Blog

		if !(helper.BindJSON(ctx, &blogData)) {
			return
		}

		_, response, err := s.Services.CreateBlog(&blogData)

		if !(helper.InternalServerError(ctx, err)) {
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"response": response})
	}
}

func (s *BlogController) UpdateBlog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var updateBlog model.Blog

		if !(helper.BindJSON(ctx, &updateBlog)) {
			return
		}

		_, response, err := s.Services.UpdateBlog(&updateBlog)

		if !(helper.InternalServerError(ctx, err)) {
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"response": response})
	}
}

func (s *DeleteBlogController) DeleteBlog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type deleteBlog struct {
			Id       string `json:"_id" binding:"required"`
			AuthorId string `json:"authorId" binding:"required"`
		}

		var blogData deleteBlog
		if !(helper.BindJSON(ctx, &blogData)) {
			return
		}

		// Call the service with the ID and AuthorId directly
		_, response, err := s.Services.DeleteBlog(blogData.Id, blogData.AuthorId)
		if !(helper.InternalServerError(ctx, err)) {
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"response": response})
	}
}

func (s *GetAllBlogController) FetchAllBlog() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Call the service with the ID and AuthorId directly
		_, response, err := s.Services.FetchAllBlog(ctx)
		if !(helper.InternalServerError(ctx, err)) {
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"response": response})
	}
}

func (s *GetOneBlogController) FetchOneBlog() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		Id := ctx.Param("blogID")

		_, response, err := s.Services.FetchOneBlog(ctx, Id)
		if !(helper.InternalServerError(ctx, err)) {
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"response": response})
	}
}

func (s *GetWriterBlogController) FetchWriterBlog() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		AuthorId := ctx.Param("authorID")

		// Call the service with the ID and AuthorId directly
		_, response, err := s.Services.FetchWriterBlog(ctx, AuthorId)
		if !(helper.InternalServerError(ctx, err)) {
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"response": response})
	}
}

func (s *SearchBlogByTitleController) SearchBlogByTitle() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		type searchBlog struct {
			Title string `json:"title" binding:"required"`
		}

		var searchData searchBlog

		if !(helper.BindJSON(ctx, &searchData)) {
			return
		}
		// Call the service with the ID and AuthorId directly
		_, response, err := s.Services.SearchBlogByTitile(ctx, searchData.Title)
		if !(helper.InternalServerError(ctx, err)) {
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"response": response})
	}
}

func (s *SearchBlogByTagController) SearchBlogByTag() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		type searchBlog struct {
			Tag string `json:"tag" binding:"required"`
		}

		var searchData searchBlog

		if !(helper.BindJSON(ctx, &searchData)) {
			return
		}
		// Call the service with the ID and AuthorId directly
		_, response, err := s.Services.SearchBlogByTag(ctx, searchData.Tag)
		if !(helper.InternalServerError(ctx, err)) {
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"response": response})
	}
}

func (s *GetAllTagController) GetAllTags() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, response, err := s.Services.GetAllTags(ctx)
		if !(helper.InternalServerError(ctx, err)) {
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"response": response})
	}
}
