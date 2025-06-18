package main

import (
	"moodmap/routes"
	"moodmap/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.InitDB()

	r := gin.Default()
	r.Static("/static", "./frontend")

	// Public routes
	r.POST("/register", routes.Register)
	r.POST("/login", routes.Login)

	// Protected routes group
	auth := r.Group("/")
	auth.Use(utils.JWTAuthMiddleware()) // ✅ JWT protection starts here

	r.POST("/mood", utils.JWTAuthMiddleware(), routes.PostMood)
	auth.GET("/all_moods", utils.JWTAuthMiddleware(), routes.GetAllMoods)

	// Optional: redirect root to frontend home
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/static/index.html")
	})

	r.Run(":8080")
}
