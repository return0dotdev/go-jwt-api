package auth

import (
	"fmt"
	"net/http"
	"os"
	"return0/jwt-api/orm"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
}

type LoginBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var json RegisterBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check user exists
	var userExist orm.User
	orm.Db.Where("username = ?", json.Username).First((&userExist))
	if userExist.ID > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User Exist!",
		})
		return
	}

	encryptPassword, _ := bcrypt.GenerateFromPassword([]byte(json.Password), 10)
	user := orm.User{Username: json.Username, Password: string(encryptPassword), Fullname: json.Fullname, Avatar: json.Avatar}

	orm.Db.Create(&user)

	if user.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "User Create Success!",
			"userId":  user.ID,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User Create failed!",
		})
	}
}

var hmacSampleSecret []byte

func Login(c *gin.Context) {
	var json LoginBody
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userExist orm.User
	orm.Db.Where("username = ?", json.Username).First((&userExist))
	if userExist.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User Does not Exist!",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(userExist.Password), []byte(json.Password))
	if err == nil {
		hmacSampleSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": userExist.ID,
			"exp":    time.Now().Add(time.Minute * 1).Unix(),
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(hmacSampleSecret)

		fmt.Println(tokenString, err)

		c.JSON(http.StatusOK, gin.H{
			"status":       "ok",
			"message":      "Login Success!",
			"access_token": tokenString,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Login failed!",
		})
		return
	}
}
