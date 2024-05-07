package handle

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/HooEP01/chat-bot/types"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/sashabaranov/go-openai"
)

type MessageHandler struct{}

type RequestData struct {
	Id       string `json:"id"`
	Question string `json:"question"`
}

func (h MessageHandler) HandleNewMessage(c *gin.Context) {
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
		respond, err = sendChatGPTRequest(data.Question)
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
}

func (h MessageHandler) HandleGeneralMessage(c *gin.Context) {
	paramType := c.Query("type")

	// Or better get from database
	questionAnswers := types.GetQuestionAnswer(paramType)

	// Generate respondData
	respondData := map[string]interface{}{
		"items":     questionAnswers,
		"timestamp": time.Now().UnixMilli(),
	}

	c.JSON(http.StatusCreated, respondData)
}

func (h MessageHandler) HandleMessageType(c *gin.Context) {
	messageType := [4]string{"sales", "partner", "support", ""}

	// Generate respondData
	respondData := map[string]interface{}{
		"items":     messageType,
		"timestamp": time.Now().UnixMilli(),
	}

	c.JSON(http.StatusCreated, respondData)
}

func notifyAdminSupport(id string) {
	// tell admin about it
}

func respondToQuestion(question string) (string, bool) {
	rules := map[string]string{
		"hello":        "Hi there! How can I help you?",
		"how are you?": "I'm just a bot, but thank for asking!",
		"goodbye":      "Goodbye! if you have more questions, feel free to ask me!",
	}

	for rule, response := range rules {
		if fuzzy.Match(rule, question) {
			return response, false
		}
	}

	return "Please wait for a minutes, a customer service representation will help you for your question.", true
}

func sendChatGPTRequest(question string) (string, error) {
	apiKey := os.Getenv("CHATGPT_API_KEY")
	if apiKey != "" {
		client := openai.NewClient(apiKey)
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: question,
					},
				},
			},
		)

		if err != nil {
			return "error", fmt.Errorf("error")
		}

		return resp.Choices[0].Message.Content, nil
	}

	return "Api key not found", fmt.Errorf("error")
}
