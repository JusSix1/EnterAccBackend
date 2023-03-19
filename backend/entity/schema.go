package entity

import (
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

// ----------User----------
type Gender struct {
	gorm.Model
	Gender string
	User   []User `gorm:"foreignKey:Gender_ID"`
}

type User struct {
	gorm.Model
	Email           string    `gorm:"uniqueIndex" valid:"email~Invalid Email format,required~Email is blank"`
	FirstName       string    `valid:"required~First name is blank"`
	LastName        string    `valid:"required~Last name is blank"`
	Password        string    `valid:"minstringlength(8)~Password must be longer than 8 characters,required~Password is blank"`
	Profile_Name    string    `valid:"maxstringlength(50)~Must be no more than 50 characters long,required~Profile name is blank"`
	Profile_Picture string    `valid:"image_valid~Please change the picture"`
	Birthday        time.Time `valid:"NotFutureTime~The day must not be the future,MoreThan18YearsAgo~You must be over 18 years old"`
	Phone_number    string    `valid:"required~Phone number is blank,matches([0-9]{10})~Phone number invalid format"`
	Gender_ID       *uint     `valid:"-"`
	Gender          Gender    `gorm:"references:id" valid:"-"`
	Account         []Account `gorm:"foreignKey:User_ID"`
	Order           []Order   `gorm:"foreignKey:User_ID"`
}

type Account_Status struct {
	gorm.Model
	Status  string
	Account []Account `gorm:"foreignKey:Account_Status_ID"`
}

type Account struct {
	gorm.Model
	ID_Account        uint           `valid:"-"`
	User_ID           *uint          `valid:"-"`
	User              User           `gorm:"references:id" valid:"-"`
	Twitter_Account   string         `valid:"-"`
	Twitter_Password  string         `valid:"-"`
	Email             string         `valid:"-"`
	Email_Password    string         `valid:"-"`
	Phone_Number      string         `valid:"-"`
	Years             uint           `valid:"-"`
	Account_Status_ID *uint          `valid:"-"`
	Account_Status    Account_Status `gorm:"references:id" valid:"-"`
	Order_ID          *uint          `valid:"-"`
	Order             Order          `gorm:"references:id" valid:"-"`
}

type Order struct {
	gorm.Model
	User_ID *uint     `valid:"-"`
	User    User      `gorm:"references:id" valid:"-"`
	Account []Account `gorm:"foreignKey:Order_ID"`
}

func init() {
	govalidator.CustomTypeTagMap.Set("DelayNow10Min", func(i interface{}, context interface{}) bool {
		t := i.(time.Time)
		return t.After(time.Now().Add(time.Minute * -10))
	})

	govalidator.CustomTypeTagMap.Set("NotFutureTime", func(i interface{}, context interface{}) bool {
		t := i.(time.Time)
		return !t.After(time.Now())
	})

	govalidator.CustomTypeTagMap.Set("MoreThan18YearsAgo", func(i interface{}, context interface{}) bool {
		t := i.(time.Time)
		ageLimit := time.Now().AddDate(-18, 0, 0)
		return t.Before(ageLimit)
	})

	govalidator.TagMap["image_valid"] = govalidator.Validator(func(str string) bool {
		return govalidator.Matches(str, "^(data:image(.+);base64,.+)$")
	})
}
