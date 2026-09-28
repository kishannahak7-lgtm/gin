package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Sociallink struct {
	Platform string `json:"platform" binding:"required,oneof=github linkedin twitter"`
	URL      string `json:"url" binding:"required,url" `
}
type CreateProfileRequest struct {
	Fullname string       `json:"fullname" binding:"required,min=3"`
	Links    []Sociallink `json:"links" binding:"required,min=1,dive"`
}

func getCustomMessage(fe validator.FieldError) string {
	switch fe.Tag() {

	case "required":
		return "this field is required"
	case "min":
		return "minimum 3 character required"
	case "oneof":
		return "platform must be one of :-linkedin,twitter ,github"
	case "url":
		return "must be a valied url"
	default:
		return "invalid value"
	}
}
func formatvalidationerror(err error) map[string]string {
	errMap := make(map[string]string)

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			errMap[fe.Namespace()] = getCustomMessage(fe)
		}
	}
	return errMap
}

func main() {
	r := gin.Default()

	r.POST("/profiles", func(c *gin.Context) {
		var req CreateProfileRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  formatvalidationerror(err),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":  "profile created successfully",
			"Fullname": req.Fullname,
			"links":    len(req.Links),
		})
	})
	r.Run(":8080")
}
