package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/shubham09anand/blogManagement/Controller"
	services "github.com/shubham09anand/blogManagement/Services"
)

type Routes struct{}

func (r *Routes) RoutesFunc(router *gin.Engine) {
	// Create instances of services
	userServices := &services.UserServices{}
	UserProfileServices := &services.UserProfileServices{}
	BlogServices := &services.BlogServices{}
	CommentsServices := &services.CommentServices{}

	// sign up user
	signupController := &controller.SignupController{
		Services: userServices,
	}

	// login up user
	loginController := &controller.LoginController{
		Services: userServices,
	}

	// setting
	settingController := &controller.SettingController{
		Services: userServices,
	}

	// get All user (writer profile)
	fetchAllUserController := &controller.FetchAllUserController{
		Services: userServices,
	}

	// profile
	profileController := &controller.ProfileController{
		Services: UserProfileServices,
	}

	// update profile
	updateProfileController := &controller.ProfileController{
		Services: UserProfileServices,
	}

	// create Blog
	blogController := &controller.BlogController{
		Services: BlogServices,
	}

	// update Blog
	updateBlogController := &controller.BlogController{
		Services: BlogServices,
	}

	// delete Blog
	deleteBlogController := &controller.DeleteBlogController{
		Services: BlogServices,
	}

	// get all Blog
	fetchAllBlogController := &controller.GetAllBlogController{
		Services: BlogServices,
	}

	// get one Blog
	fetchOneBlogController := &controller.GetOneBlogController{
		Services: BlogServices,
	}

	// get a specfic writer Blog
	fetchWriterBlogController := &controller.GetWriterBlogController{
		Services: BlogServices,
	}

	// search blog by Title
	searchBlogByTitleController := &controller.SearchBlogByTitleController{
		Services: BlogServices,
	}

	// search blog by Tag
	searchBlogByTagController := &controller.SearchBlogByTagController{
		Services: BlogServices,
	}

	// make comment
	makeCommentController := &controller.MakeCommentController{
		Services: CommentsServices,
	}

	// delete comment
	deleteCommentController := &controller.DeleteCommentController{
		Services: CommentsServices,
	}

	// get all tags for search
	getAllTasController := &controller.GetAllTagController{
		Services: BlogServices,
	}

	// get writer name for search tags
	getAllWriterNameController := &controller.FetchAllWriterNameController{
		Services: userServices,
	}

	getCommentsController := &controller.GetCommentController{
		Services: CommentsServices,
	}

	// this is for getting requested(only one) writer profile details
	getUserProfileController := &controller.FetchUserProfileController{
		Services: userServices,
	}

	// to get the user photo for header
	getUserPhotoController := &controller.GetUerPhotoController{
		Services: UserProfileServices,
	}

	// Define routes with correct handler functions
	router.POST("/signup", signupController.Signup())
	router.POST("/login", loginController.Login())
	router.POST("/profile", profileController.Profile())
	router.POST("/updateProfile", updateProfileController.UpdateProfile())

	router.POST("/createBlog", blogController.CreateBlog())
	router.POST("/updateBlog", updateBlogController.UpdateBlog())
	router.POST("/deleteBlog", deleteBlogController.DeleteBlog())

	router.POST("/searchBlogByTitle", searchBlogByTitleController.SearchBlogByTitle())
	router.POST("/searchBlogByTag", searchBlogByTagController.SearchBlogByTag())

	router.POST("/setting", settingController.Setting())

	router.POST("/getWriters", fetchAllUserController.FetchAllUser())

	router.POST("/makeComments", makeCommentController.MakeComment())
	router.POST("/deleteComments", deleteCommentController.DeleteComment())

	router.POST("/getAllBlogs", fetchAllBlogController.FetchAllBlog())
	router.GET("/getOneBlog/:blogID", fetchOneBlogController.FetchOneBlog())
	router.GET("/getWriterBlog/:authorID", fetchWriterBlogController.FetchWriterBlog())

	router.POST("/getAllTags", getAllTasController.GetAllTags())
	router.POST("/getAllWriterName", getAllWriterNameController.FetchAllWriterName())

	router.GET("/getComments/:blogId", getCommentsController.GetComments())

	router.GET("/getWriterProfile/:userId", getUserProfileController.GetWriterProfile())

	router.POST("/getUserPhoto", getUserPhotoController.GetUserPhoto())
}
