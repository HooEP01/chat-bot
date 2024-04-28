package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestData struct {
	Id       string `json:"id"`
	Question string `json:"question"`
}

func main() {
	router := gin.Default()

	router.POST("/message", func(c *gin.Context) {
		var data RequestData
		err := c.BindJSON(&data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Process data

		// Save to database

		c.JSON(http.StatusCreated, data)
	})

	router.Run(":8080")
}
