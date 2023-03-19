package main

import (
	account_controller "github.com/JusSix1/TwitterAccountDataBase/controller/account"
	login_controller "github.com/JusSix1/TwitterAccountDataBase/controller/login"
	order_controller "github.com/JusSix1/TwitterAccountDataBase/controller/order"
	user_controller "github.com/JusSix1/TwitterAccountDataBase/controller/user"
	"github.com/JusSix1/TwitterAccountDataBase/entity"
	"github.com/JusSix1/TwitterAccountDataBase/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {

	entity.SetupDatabase()

	r := gin.Default()
	r.Use(CORSMiddleware())

	// login User Route
	r.POST("/login/user", login_controller.LoginUser)
	r.POST("/users", user_controller.CreateUser)
	r.GET("/genders", user_controller.ListGenders)

	// // login Admin Route
	// r.POST("/login/admin", login_controller.LoginAdmin)

	router := r.Group("/")
	{
		protected := router.Use(middlewares.Authorizes())
		{
			protected.GET("/user/:email", user_controller.GetUser)
			protected.GET("/usersprofilepicture/:email", user_controller.GetUserProfilePicture)
			protected.PATCH("/users", user_controller.UpdateUser)
			protected.PATCH("/usersPassword", user_controller.UpdateUserPassword)
			protected.DELETE("/users/:email", user_controller.DeleteUser)

			protected.POST("/account", account_controller.CreateAccount)
			protected.GET("/all-account/:email", account_controller.GetAllAccount)
			protected.GET("/unsold-account/:email", account_controller.GetUnsoldAccount)
			protected.GET("/account-in-order/:id", account_controller.GetAccountInOrder)
			protected.DELETE("/account", account_controller.DeleteAccount)

			protected.POST("/order/:email", order_controller.CreateOrder)
			protected.GET("/order/:email", order_controller.GetOrder)
		}
	}

	// Run the server
	r.Run()

}

func CORSMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")

		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {

			c.AbortWithStatus(204)

			return

		}

		c.Next()

	}

}
