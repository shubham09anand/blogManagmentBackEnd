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
		var updateprofileData model.UserProfile

		if !(helper.BindJSON(ctx, &updateprofileData)) {
			return
		}

		_, response, err := s.Services.UpdateProfile(&updateprofileData)

		if !(helper.InternalServerError(ctx, err)) {
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"response": response})
	}
}
