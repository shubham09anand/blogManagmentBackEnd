package main

import (
	"github.com/gin-gonic/gin"
	routes "github.com/shubham09anand/blogManagement/Routes"
)

func main() {

	r := gin.Default()

	port := "4000"

	// Initialize routes
	r.Use(CORSMiddleware())
	routes := routes.Routes{}
	routes.RoutesFunc(r)

	// Start the server
	r.Run(":" + port)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
