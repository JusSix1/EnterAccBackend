package controller

import (
	"net/http"

	"github.com/JusSix1/TwitterAccountDataBase/entity"
	"github.com/asaskevich/govalidator"
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

	var payload struct {
		AccountToOrder []AccountToOrder `json:"accountToOrder"`
		DataSlip       struct {
			Slip string `json:"Slip"`
		} `json:"dataSlip"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
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
		Slip:    payload.DataSlip.Slip,
	}

	if err := entity.DB().Create(&newOrder).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := entity.DB().Raw("SELECT * FROM orders WHERE user_id = ? ORDER BY id DESC LIMIT 1", user.ID).Find(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i := 0; i < len(payload.AccountToOrder); i++ {

		if err := entity.DB().Raw("SELECT * FROM account_statuses WHERE status = 'Sold'").Find(&accoountStatus).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// create new object for create new record
		updateAccount := entity.Account{
			Account_Status_ID: &accoountStatus.ID,
			Order_ID:          &order.ID,
		}

		if err := entity.DB().Where("id = ?", payload.AccountToOrder[i].Account_ID).Updates(&updateAccount).Error; err != nil {
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

// PATCH /order
func UpdateOrder(c *gin.Context) {
	var order entity.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// update user fields that are allowed to be updated
	updateOrder := entity.Order{
		Slip: order.Slip,
	}

	// validate user
	if _, err := govalidator.ValidateStruct(updateOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := entity.DB().Where("id = ?", order.ID).Updates(&updateOrder).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

// GET /all-order-admin
func GetOrderAdmin(c *gin.Context) {
	var order []entity.Order

	if err := entity.DB().Raw("SELECT * FROM orders").Find(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}
