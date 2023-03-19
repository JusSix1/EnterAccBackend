package controller

import (
	"net/http"

	"github.com/JusSix1/TwitterAccountDataBase/entity"
	"github.com/gin-gonic/gin"
)

/* --- ระบบ Login ---*/
// GET /login
func ListUsersLogin(c *gin.Context) {
	var user []entity.User

	if err := entity.DB().Raw("SELECT * FROM users").Scan(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
