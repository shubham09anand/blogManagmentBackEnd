package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shubham09anand/blogManagement/routes"
)

func main() {

	r := gin.Default()

	// Initialize routes
	r.Use(CORSMiddleware())
	routes := routes.Routes{}
	routes.RoutesFunc(r)

	// Start the server
	r.Run(":4000")
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
