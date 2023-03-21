package controller

import (
	"net/http"

	"github.com/JusSix1/TwitterAccountDataBase/entity"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// POST /users
func CreateAdmin(c *gin.Context) {
	var admin entity.Admin
	var adminNamecheck entity.Admin

	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if tx := entity.DB().Where("admin_name = ?", admin.Admin_Name).First(&adminNamecheck); !(tx.RowsAffected == 0) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This Admin Name has already been taken."})
		return
	}

	// create new object for create new record
	newAdmin := entity.Admin{
		Admin_Name: admin.Admin_Name,
		Password:   admin.Password,
		Big:        admin.Big,
	}

	// validate user
	if _, err := govalidator.ValidateStruct(newAdmin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hashing after validate
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newAdmin.Password), 12)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error hashing password"})
		return
	}

	newAdmin.Password = string(hashPassword)

	if err := entity.DB().Create(&newAdmin).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": admin})

}

// GET /admin-list
func GetAdminList(c *gin.Context) {
	var admin []entity.Admin

	if err := entity.DB().Raw("SELECT * FROM admins").Find(&admin).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": admin})
}

// PATCH /admin-right/:adminname
func UpdateBigAdmin(c *gin.Context) {
	var admin entity.Admin

	adminName := c.Param("adminname")

	if tx := entity.DB().Where("admin_name = ?", adminName).First(&admin); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Admin not found"})
		return
	}

	if admin.Big {

		type UpdateAdmin struct {
			ID uint
		}

		var updateAdmin []UpdateAdmin

		if err := c.ShouldBindJSON(&updateAdmin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updateBigAdmin := entity.Admin{
			Big: true,
		}

		for i := 0; i < len(updateAdmin); i++ {

			if int(admin.ID) != int(updateAdmin[i].ID) {

				if err := entity.DB().Where("id = ?", updateAdmin[i].ID).Updates(&updateBigAdmin).Error; err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

			}
		}

		c.JSON(http.StatusOK, gin.H{"data": updateAdmin})

	} else {

		c.JSON(http.StatusBadRequest, gin.H{"error": "You have no rights"})

	}

}

// DELETE /admin/:adminname
func DeleteAdmin(c *gin.Context) {
	var admin entity.Admin

	adminName := c.Param("adminname")

	if tx := entity.DB().Where("admin_name = ?", adminName).First(&admin); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Admin not found"})
		return
	}

	if admin.Big {

		type DeleteAdmin struct {
			ID uint
		}

		var deleteAdmin []DeleteAdmin

		if err := c.ShouldBindJSON(&deleteAdmin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for i := 0; i < len(deleteAdmin); i++ {

			if int(admin.ID) != int(deleteAdmin[i].ID) {

				if tx := entity.DB().Exec("DELETE FROM admins WHERE id = ?", deleteAdmin[i].ID); tx.RowsAffected == 0 {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Admin not found"})
					return
				}

			}
		}

		c.JSON(http.StatusOK, gin.H{"data": deleteAdmin})

	} else {

		c.JSON(http.StatusBadRequest, gin.H{"error": "You have no rights"})

	}

}
