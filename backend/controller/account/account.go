package controller

import (
	"net/http"

	"github.com/JusSix1/TwitterAccountDataBase/entity"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// POST /account
func CreateAccount(c *gin.Context) {
	var user entity.User
	var accountStatus entity.Account_Status
	var account entity.Account

	type AccountImport struct {
		Email_User        string
		Twitter_Account   string
		Twitter_Password  string
		Email_Accont      string
		Email_Password    string
		Phone_Number      string
		Years             uint `gorm:"required~Year is blank"`
		Account_Status_ID uint
	}

	var accountImport []AccountImport

	if err := c.ShouldBindJSON(&accountImport); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i := 0; i < len(accountImport); i++ {

		if tx := entity.DB().Where("email = ?", accountImport[i].Email_User).First(&user); tx.RowsAffected == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}

		if tx := entity.DB().Where("id = ?", accountImport[i].Account_Status_ID).First(&accountStatus); tx.RowsAffected == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Account status not found"})
			return
		}

		if _, err := govalidator.ValidateStruct(accountImport[i]); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := entity.DB().Raw("SELECT * FROM accounts WHERE user_id = ? ORDER BY id DESC LIMIT 1", user.ID).Find(&account).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// create new object for create new record
		newAccount := entity.Account{
			User_ID:           &user.ID,
			ID_Account:        account.ID_Account + 1,
			Twitter_Account:   accountImport[i].Twitter_Account,
			Twitter_Password:  accountImport[i].Twitter_Password,
			Email:             accountImport[i].Email_Accont,
			Email_Password:    accountImport[i].Email_Password,
			Phone_Number:      accountImport[i].Phone_Number,
			Years:             accountImport[i].Years,
			Account_Status_ID: &accountStatus.ID,
		}

		if err := entity.DB().Create(&newAccount).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": accountImport})

}

// GET /all-account/:email
func GetAllAccount(c *gin.Context) {
	var user entity.User
	var account []entity.Account

	email := c.Param("email")

	if tx := entity.DB().Where("email = ?", email).First(&user); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	if err := entity.DB().Preload("Account_Status").Raw("SELECT * FROM accounts WHERE user_id = ? ORDER BY id_account DESC", user.ID).Find(&account).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": account})
}

// GET /unsold-account/:email
func GetUnsoldAccount(c *gin.Context) {
	var user entity.User
	var account []entity.Account
	var accoountStatus entity.Account_Status

	email := c.Param("email")

	if tx := entity.DB().Where("email = ?", email).First(&user); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	if err := entity.DB().Raw("SELECT * FROM account_statuses WHERE status = 'Unsold'").Find(&accoountStatus).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := entity.DB().Preload("Account_Status").Raw("SELECT * FROM accounts WHERE user_id = ? AND account_status_id = ? ORDER BY id_account DESC", user.ID, accoountStatus.ID).Find(&account).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": account})
}

// GET /account-in-order
func GetAccountInOrder(c *gin.Context) {
	var account []entity.Account

	id := c.Param("id")

	if err := entity.DB().Raw("SELECT * FROM accounts WHERE order_id = ? ORDER BY id_account DESC", id).Find(&account).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": account})
}

// DELETE /account
func DeleteAccount(c *gin.Context) {

	type DeleteAccount struct {
		ID uint
	}

	var deleteAccount []DeleteAccount

	if err := c.ShouldBindJSON(&deleteAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i := 0; i < len(deleteAccount); i++ {

		if tx := entity.DB().Exec("DELETE FROM accounts WHERE id = ?", deleteAccount[i].ID); tx.RowsAffected == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "basket not found"})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{"data": deleteAccount})
}
