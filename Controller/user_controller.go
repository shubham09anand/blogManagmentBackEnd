package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	helper "github.com/shubham09anand/blogManagement/Helper"
	model "github.com/shubham09anand/blogManagement/Model"
	services "github.com/shubham09anand/blogManagement/Services"
)

type SignupController struct {
	Services *services.UserServices
}

type LoginController struct {
	Services *services.UserServices
}

type SettingController struct {
	Services *services.UserServices
}

type FetchAllUserController struct {
	Services *services.UserServices
}

// this is for getting all writer name for sarch
type FetchAllWriterNameController struct {
	Services *services.UserServices
}

// this is for getting requested(only one) writer profile details
type FetchUserProfileController struct {
	Services *services.UserServices
}

func (s *SignupController) Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var signupData *model.UserSignup

		if !(helper.BindJSON(ctx, &signupData)) {
			return
		}

		_, response, err := s.Services.Signup(signupData)

		if !(helper.InternalServerError(ctx, err)) {
			return
		}

		// fmt.Println("Response:", response)

		ctx.JSON(http.StatusOK, gin.H{"response": response})
	}
}

func (s *LoginController) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type loginData struct {
			UserName string `json:"userName" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		var data loginData

		if !(helper.BindJSON(ctx, &data)) {
			return
		}

		// Call the Login method instead of Signup
		_, response, err := s.Services.Login(data.UserName, data.Password)
		if !(helper.InternalServerError(ctx, err)) {
			return
		}

		// fmt.Println("Response:", response)

		ctx.JSON(http.StatusOK, gin.H{"response": response})
	}
}

func (s *SettingController) Setting() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type settingData struct {
			Id       string `json:"_id" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		var data settingData

		if !(helper.BindJSON(ctx, &data)) {
			return
		}

		// Call the Login method instead of Signup
		_, response, err := s.Services.Setting(data.Id, data.Password)
		if !(helper.InternalServerError(ctx, err)) {
			return
		}

		// fmt.Println("Response:", response)

		ctx.JSON(http.StatusOK, gin.H{"response": response})
	}
}

func (s *FetchAllUserController) FetchAllUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Call the Login method instead of Signup
		_, response, err := s.Services.FetchAllUser(ctx)
		if !(helper.InternalServerError(ctx, err)) {
			return
		}

		// fmt.Println("Response:", response)

		ctx.JSON(http.StatusOK, gin.H{"response": response})
	}
}

func (s *FetchAllWriterNameController) FetchAllWriterName() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Call the Login method instead of Signup
		_, response, err := s.Services.GetAllWriterName(ctx)
		if !(helper.InternalServerError(ctx, err)) {
			return
		}

		// fmt.Println("Response:", response)

		ctx.JSON(http.StatusOK, gin.H{"response": response})
	}
}

// this is for getting requested(only one) writer profile details
func (s *FetchUserProfileController) GetWriterProfile() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		id := ctx.Param("userId")

		_, response, err := s.Services.GetWriterProfile(ctx, id)
		if !(helper.InternalServerError(ctx, err)) {
			return
		}

		// fmt.Println("Response:", response)

		ctx.JSON(http.StatusOK, gin.H{"response": response})
	}
}
