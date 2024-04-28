package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

type RequestData struct {
	Id       string `json:"id"`
	Question string `json:"question"`
}

var loc, _ = time.LoadLocation("Asia/Kuala_Lumpur")

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
		respond, notification := respondToQuestion(data.Question)

		// Notify admin support
		if notification {
			notifyAdminSupport(data.Id)
		}

		// Save to database

		// Generate respondData
		respondData := map[string]interface{}{
			"id":         "uuid", // generate from database
			"request_id": data.Id,
			"answer":     respond,
			"timestamp":  time.Now().UnixMilli(),
		}

		c.JSON(http.StatusCreated, respondData)
	})

	router.Run(":8080")
}

func respondToQuestion(question string) (string, bool) {
	rules := map[string]string{
		"hello":        "Hi there! How can I help you?",
		"how are you?": "I'm just a bot, but thank for asking!",
		"goodbye":      "Goodbye! if you have more questions, feel free to ask me!",
	}

	// Spelling and fuzzy check

	for rule, response := range rules {
		if fuzzy.Match(rule, question) {
			return response, false
		}
	}

	return "Please wait for a minutes, a customer service representation will help you for your question.", true
}

func notifyAdminSupport(id string) {

}
