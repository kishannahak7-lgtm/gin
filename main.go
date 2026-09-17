package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderRequest struct {
	ProductName string  `json:"product_name" binding:"required"`
	Quantity    int     `json:"quantity" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
}
type Car struct {
	ID    int     `json:"id"`
	Brand string  `json:"brand" binding:"required"`
	Model string  `json:"model" binding:"required"`
	Price float64 `json:"price" binding:"required"`
}

var cars = []Car{
	{ID: 1, Brand: "Toyota", Model: "fortuner", Price: 25000},
	{ID: 2, Brand: "Honda", Model: "Civic", Price: 22000},
	{ID: 3, Brand: "Ford", Model: "Mustang", Price: 35000},
}

type Movie struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Rating float64 `json:"rating"`
}

var movies = []Movie{
	{ID: 1, Title: "Movie 1", Rating: 8.5},
	{ID: 2, Title: "Movie 2", Rating: 7.8},
	{ID: 3, Title: "Movie 3", Rating: 9.2},
}

func main() {
	r := gin.Default()

	r.GET("/hotel/:city", func(c *gin.Context) {
		city := c.Param("city")
		rating := c.DefaultQuery("rating", "any")
		price := c.DefaultQuery("price", "0")
		wifi := c.DefaultQuery("wifi", "false")

		c.JSON(http.StatusOK, gin.H{
			"city":   city,
			"rating": rating,
			"price":  price,
			"wifi":   wifi,
		})
	})

	r.POST("/order", func(c *gin.Context) {
		var req OrderRequest

		// Read raw body to provide clearer error messages when parsing fails
		raw, _ := io.ReadAll(c.Request.Body)
		// restore Body so gin/bind can read it again if needed
		c.Request.Body = io.NopCloser(io.MultiReader(bytes.NewReader(raw)))

		if err := json.Unmarshal(raw, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request",
				"details": err.Error(),
				"body":    string(raw),
			})
			return
		}
		if req.Quantity <= 0 || req.Price <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Quantity & price it must > 0",
			})
			return
		}
		totalAmount := float64(req.Quantity) * req.Price

		c.JSON(http.StatusCreated, gin.H{
			"message":      "Order created successfully",
			"product_name": req.ProductName,
			"quantity":     req.Quantity,
			"price":        req.Price,
			"total_amount": totalAmount,
		})
	})

	r.PUT("/car/:id", func(c *gin.Context) {
		idStr := c.Param("id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "id must be a number",
			})
			return
		}

		var req Car
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "id, brand, model, price is required",
				"details": err.Error(),
			})
			return
		}

		for i, car := range cars {
			if car.ID == id {
				req.ID = id
				cars[i] = req

				c.JSON(http.StatusOK, gin.H{
					"message": "car updated successfully",
					"car":     req,
				})
				return
			}

		}
		c.JSON(http.StatusNotFound, gin.H{
			"error": "car not found",
		})

	})

	r.PATCH("/movie/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "id must be a number",
			})
			return
		}
		var input struct {
			Title  string  `json:"title"`
			Rating float64 `json:"rating"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalied json payload",
			})
			return
		}

		for i, movie := range movies {
			if movie.ID == id {
				if input.Title != "" {
					movies[i].Title = input.Title
				}
				if input.Rating != 0 {
					movies[i].Rating = input.Rating
				}
				c.JSON(http.StatusOK, gin.H{
					"message": "Movie updated successfully",
					"movie":   movies[i],
				})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Movie not found",
		})

	})

	r.DELETE("/movie/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "id must be a number",
			})
			return
		}

		for i, movie := range movies {
			if movie.ID == id {
				movies = append(movies[:i], movies[i+1:]...)
				c.JSON(http.StatusOK, gin.H{
					"message":       "Movie deleted successfully",
					"deleted_movie": id,
				})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Movie not found",
		})

	})
	r.Run(":8080")

}
