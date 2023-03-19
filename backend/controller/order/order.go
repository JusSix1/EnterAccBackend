package controller

import (
	"net/http"

	"github.com/JusSix1/TwitterAccountDataBase/entity"
	"github.com/gin-gonic/gin"
)

// POST /order/:email
func CreateOrder(c *gin.Context) {
	var user entity.User
	var order entity.Order
	var accoountStatus entity.Account_Status

	email := c.Param("email")

	type AccountToOrder struct {
		Account_ID uint
	}

	var accountToOrder []AccountToOrder

	if err := c.ShouldBindJSON(&accountToOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if tx := entity.DB().Where("email = ?", email).First(&user); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// create new object for create new record
	newOrder := entity.Order{
		User_ID: &user.ID,
	}

	if err := entity.DB().Create(&newOrder).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := entity.DB().Raw("SELECT * FROM orders WHERE user_id = ? ORDER BY id DESC LIMIT 1", user.ID).Find(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i := 0; i < len(accountToOrder); i++ {

		if err := entity.DB().Raw("SELECT * FROM account_statuses WHERE status = 'Sold'").Find(&accoountStatus).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// create new object for create new record
		updateAccount := entity.Account{
			Account_Status_ID: &accoountStatus.ID,
			Order_ID:          &order.ID,
		}

		if err := entity.DB().Where("id = ?", accountToOrder[i].Account_ID).Updates(&updateAccount).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

// GET /order/:email
func GetOrder(c *gin.Context) {
	var user entity.User
	var order []entity.Order

	email := c.Param("email")

	if tx := entity.DB().Where("email = ?", email).First(&user); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	if err := entity.DB().Raw("SELECT * FROM orders WHERE user_id = ? ORDER BY id DESC", user.ID).Find(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}
