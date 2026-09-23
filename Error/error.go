package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CreateProductInput struct {
	Title    string  `json:"title" binding:"required,min=5"`
	Price    float64 `json:"price" binding:"required,gt=0"`
	Category string  `json:"category" binding:"required,oneof=electronics clothing books"`
}

func FormatValidationError(err error) map[string]string {
	errorMap := make(map[string]string)

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			errorMap[fe.Field()] = GetCustomMessage(fe)
		}
	}
	return errorMap
}

func GetCustomMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "gt":
		return "price must be greater than 0"
	case "oneof":
		return "Category must be one of: electronics, clothing, books"
	default:
		return "Invalid value"
	}
}

func main() {
	r := gin.Default()

	r.POST("/product", func(c *gin.Context) {
		var input CreateProductInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "Error",
				"error":  FormatValidationError(err),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "product created successfully",
		})
	})

	r.Run(":8080")
}
