package controller

import (
	"net/http"
	"time"

	"github.com/JusSix1/TwitterAccountDataBase/entity"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// POST /users
func CreateUser(c *gin.Context) {
	var user entity.User
	var emailCheck entity.User
	var gender entity.Gender

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if tx := entity.DB().Where("id = ?", user.Gender_ID).First(&gender); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please select gender"})
		return
	}

	if tx := entity.DB().Where("email = ?", user.Email).First(&emailCheck); !(tx.RowsAffected == 0) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This email has already been taken."})
		return
	}

	// create new object for create new record
	newUser := entity.User{
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Password:        user.Password,
		Profile_Name:    user.Profile_Name,
		Profile_Picture: user.Profile_Picture,
		Gender:          gender,
		Birthday:        user.Birthday,
		Phone_number:    user.Phone_number,
	}

	// validate user
	if _, err := govalidator.ValidateStruct(newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hashing after validate
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 12)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error hashing password"})
		return
	}

	newUser.Password = string(hashPassword)

	if err := entity.DB().Create(&newUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})

}

// GET /user/:email
func GetUser(c *gin.Context) {
	var user entity.User
	email := c.Param("email")

	if err := entity.DB().Preload("Gender").Raw("SELECT id,email,first_name,last_name,profile_name,profile_picture,birthday,gender_id,phone_number FROM users WHERE email = ? AND deleted_at IS NULL", email).Find(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// GET /usersprofilepicture/:email
func GetUserProfilePicture(c *gin.Context) {
	var user entity.User
	email := c.Param("email")

	if err := entity.DB().Raw("SELECT profile_picture FROM users WHERE email = ? AND deleted_at IS NULL", email).Find(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// PATCH /users
func UpdateUser(c *gin.Context) {
	var user entity.User
	var gender entity.Gender

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if tx := entity.DB().Where("id = ?", user.Gender_ID).First(&gender); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gender not found"})
		return
	}

	// create new struct with fields to validate
	validateUser := struct {
		FirstName       string    `valid:"required~First name is blank"`
		LastName        string    `valid:"required~Last name is blank"`
		Profile_Name    string    `valid:"maxstringlength(50)~Must be no more than 50 characters long,required~Profile name is blank"`
		Profile_Picture string    `valid:"image_valid~Please change the picture"`
		Birthday        time.Time `valid:"NotFutureTime~The day must not be the future,MoreThan18YearsAgo~You must be over 18 years old"`
		Phone_number    string    `valid:"required~Phone number is blank,matches([0-9]{10})~Phone number invalid format"`
	}{FirstName: user.FirstName, LastName: user.LastName, Profile_Name: user.Profile_Name, Profile_Picture: user.Profile_Picture, Birthday: user.Birthday, Phone_number: user.Phone_number}

	// validate user
	if _, err := govalidator.ValidateStruct(validateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// update user fields that are allowed to be updated
	updateUser := entity.User{
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Profile_Name:    user.Profile_Name,
		Profile_Picture: user.Profile_Picture,
		Gender:          gender,
		Birthday:        user.Birthday,
		Phone_number:    user.Phone_number,
	}

	if err := entity.DB().Where("email = ?", user.Email).Updates(&updateUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// PATCH /usersPassword
func UpdateUserPassword(c *gin.Context) {

	type PasswordUser struct {
		Email       string `gorm:"uniqueIndex" valid:"email~Invalid Email format,required~Email is blank"`
		OldPassword string `valid:"minstringlength(8)~Password must be longer than 8 characters,required~Password is blank"`
		NewPassword string `valid:"minstringlength(8)~Password must be longer than 8 characters,required~Password is blank"`
	}

	var passwordUser PasswordUser
	var oldPassword entity.User

	if err := c.ShouldBindJSON(&passwordUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if tx := entity.DB().Where("email = ?", passwordUser.Email).Last(&oldPassword); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// validate PasswordUser
	if _, err := govalidator.ValidateStruct(passwordUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ตรวจสอบรหัสผ่าน
	err := bcrypt.CompareHashAndPassword([]byte(oldPassword.Password), []byte(passwordUser.OldPassword))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Old password is incorrect"})
		return
	}

	updatePasswordUser := entity.User{
		Password: passwordUser.NewPassword,
	}

	if !(passwordUser.NewPassword[0:7] == "$2a$12$") { // เช็คว่ารหัสที่ผ่านเข้ามามีการ encrypt แล้วหรือยัง หากมีการ encrypt แล้วจะไม่ทำการ encrypt ซ้ำ
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(updatePasswordUser.Password), 12)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error hashing password"})
			return
		}
		print("HASH!!!!")
		updatePasswordUser.Password = string(hashPassword)
	} else {
		print("NOT HASH!!!")
	}

	if err := entity.DB().Where("email = ?", passwordUser.Email).Updates(&updatePasswordUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": passwordUser})

}

// DELETE /users/:email
func DeleteUser(c *gin.Context) {
	email := c.Param("email")

	type PasswordUser struct {
		Password string `valid:"minstringlength(8)~Password must be longer than 8 characters,required~Password is blank"`
	}

	var passwordUser PasswordUser
	var Password entity.User

	if err := c.ShouldBindJSON(&passwordUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if tx := entity.DB().Where("email = ?", email).Last(&Password); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// validate PasswordUser
	if _, err := govalidator.ValidateStruct(passwordUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ตรวจสอบรหัสผ่าน
	err := bcrypt.CompareHashAndPassword([]byte(Password.Password), []byte(passwordUser.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Old password is incorrect"})
		return
	}

	if tx := entity.DB().Where("email = ?", email).Delete(&entity.User{}); tx.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	deleteUser := entity.User{
		Email: Password.Email + "-",
	}

	if err := entity.DB().Where("email = ?", Password.Email).Updates(&deleteUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": email})
}
