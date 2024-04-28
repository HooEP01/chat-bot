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

type QuestionAnswer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
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
		respond, isNotifyNeeded := respondToQuestion(data.Question)

		// Notify admin support
		if isNotifyNeeded {
			notifyAdminSupport(data.Id)

			// Or chatgpt api
			respond = sendChatGPTRequest(data.Question)
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

	router.GET("/message", func(c *gin.Context) {
		paramType := c.Query("type")

		// Or better get from database
		questionAnswers := getQuestionAnswers(paramType)

		// Generate respondData
		respondData := map[string]interface{}{
			"items":     questionAnswers,
			"timestamp": time.Now().UnixMilli(),
		}

		c.JSON(http.StatusCreated, respondData)
	})

	router.GET("/type", func(c *gin.Context) {
		messageType := [4]string{"sales", "partner", "support", ""}

		// Generate respondData
		respondData := map[string]interface{}{
			"items":     messageType,
			"timestamp": time.Now().UnixMilli(),
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

func sendChatGPTRequest(question string) string {
	return "This is chat gpt respond"
}

func getQuestionAnswers(paramType string) []QuestionAnswer {
	questionAnswerList := map[string][]QuestionAnswer{
		"sales": []QuestionAnswer{
			{Question: "Plan A", Answer: "Plan A is ..."},
			{Question: "Plan B", Answer: "Plan B is ..."},
		},
		"partner": []QuestionAnswer{
			{Question: "Advertisement", Answer: "To advertise in here ..."},
			{Question: "Integration", Answer: "Here is our criteria, please contact ..."},
		},
		"support": []QuestionAnswer{
			{Question: "Mobile App", Answer: "To solve there problem, please go here ..."},
			{Question: "Admin Portal", Answer: "To solve there problem, please go here ..."},
		},
		"": []QuestionAnswer{
			{Question: "default", Answer: "This is default"},
		},
	}

	questionAnswers, ok := questionAnswerList[paramType]

	if !ok {
		return []QuestionAnswer{}
	}

	return questionAnswers
}
