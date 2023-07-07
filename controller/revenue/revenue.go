package controller

import (
	"net/http"

	"github.com/JusSix1/TwitterAccountDataBase/entity"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// POST /revenue/:email
func CreateRevenue(c *gin.Context) {
	var revenue entity.Revenue
	var user entity.User

	email := c.Param("email")

	if err := c.ShouldBindJSON(&revenue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if tx := entity.DB().Where("email = ?", email).First(&user); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// create new object for create new record
	newRevenue := entity.Revenue{
		User:   user,
		Income: revenue.Income,
	}

	// validate user
	if _, err := govalidator.ValidateStruct(newRevenue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := entity.DB().Create(&newRevenue).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})

}

// GET /revenue/:email
func GetRevenue(c *gin.Context) {
	var user entity.User
	var revenue []entity.Revenue

	email := c.Param("email")

	if tx := entity.DB().Where("email = ?", email).First(&user); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	if err := entity.DB().Raw("SELECT * FROM revenues WHERE user_id = ? ORDER BY id DESC", user.ID).Find(&revenue).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": revenue})
}

// PATCH /revenue
func UpdateRevenue(c *gin.Context) {
	var revenue entity.Revenue

	if err := c.ShouldBindJSON(&revenue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// update user fields that are allowed to be updated
	updateRevenue := entity.Revenue{
		Income: revenue.Income,
	}

	// validate user
	if _, err := govalidator.ValidateStruct(updateRevenue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := entity.DB().Where("id = ?", revenue.ID).Updates(&updateRevenue).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": revenue})
}
