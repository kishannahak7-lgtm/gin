package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	message string `json:"message"`
	Name  string `json:"name" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Email string `json:"email" binding:"required,email"`
	Role string `json:"role" binding:"required,oneof=admin user"`
}

func main(){
	r := gin.Default()

	r.POST("/register",func(c *gin.Context){
		var user User
		if err := c.ShouldBindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest,gin.H{
				"error":"Validation failed",
				"details":err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK,gin.H{
			"message":"User registered successfully",
			"name":user.Name,
			"password":user.Password,
			
			"role":user.Role,
			"email":user.Email,
		})
	})
	r.Run(":8080")
}

