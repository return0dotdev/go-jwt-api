package user

import (
	"net/http"
	"return0/jwt-api/orm"

	"github.com/gin-gonic/gin"
)

func ReadAll(c *gin.Context) {
	var users []orm.User

	orm.Db.Find(&users)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "user": users})
}

func Profile(c *gin.Context) {
	userId := c.MustGet("userId").(float64)
	var user []orm.User

	orm.Db.First(&user, userId)
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "User Read Success", "user": user})
}
