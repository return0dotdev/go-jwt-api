package main

import (
	"log"
	AuthController "return0/jwt-api/controller/auth"
	UserController "return0/jwt-api/controller/user"

	"return0/jwt-api/middleware"
	"return0/jwt-api/orm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	orm.InitDb()
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/register", AuthController.Register)
	r.POST("/login", AuthController.Login)

	authorization := r.Group("/users", middleware.JWTAuth())
	authorization.GET("/readAll", UserController.ReadAll)
	authorization.GET("/profile", UserController.Profile)

	r.Run()
}
