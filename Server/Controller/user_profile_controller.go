package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	helper "github.com/shubham09anand/blogManagement/helper"
	model "github.com/shubham09anand/blogManagement/model"
	services "github.com/shubham09anand/blogManagement/services"
)

type ProfileController struct {
	Services *services.UserProfileServices
}

type UpdateProfileController struct {
	Services *services.UserProfileServices
}

type GetUerPhotoController struct {
	Services *services.UserProfileServices
}

func (s *ProfileController) Profile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var profileData model.UserProfile

		if !(helper.BindJSON(ctx, &profileData)) {
			return
		}

		_, response, err := s.Services.Profile(&profileData)

		if !(helper.InternalServerError(ctx, err)) {
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"response": response})
	}
}

func (s *ProfileController) UpdateProfile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestData struct {
			FirstName   string            `json:"firstName" binding:"required"`
			LastName    string            `json:"lastName" binding:"required"`
			UserProfile model.UserProfile `json:"userProfile" binding:"required"`
		}

		// Bind JSON payload
		if err := ctx.ShouldBindJSON(&requestData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":   400,
				"response": "Binding Error",
				"error":    err.Error(),
			})
			return
		}

		// Call the service to update the profile
		_, response, err := s.Services.UpdateProfile(requestData.FirstName, requestData.LastName, &requestData.UserProfile)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":   500,
				"response": response.Response,
				"error":    err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"response": response})
	}
}

func (s *GetUerPhotoController) GetUserPhoto() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type userIdModel struct {
			Id string `json:"userId" binding:"required"`
		}

		var userId userIdModel

		if err := ctx.ShouldBindJSON(&userId); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		photo, response, err := s.Services.GetUserPhoto(userId.Id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user photo"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"response": response, "photo": photo})
	}
}
