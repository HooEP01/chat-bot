package handle

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/HooEP01/chat-bot/types"
	"github.com/gin-gonic/gin"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/sashabaranov/go-openai"
)

type RequestData struct {
	Id       string `json:"id"`
	Question string `json:"question"`
}

type ResultData struct {
	Data string `json:"data"`
	Err  error  `json:"error,omitempty"`
}

func HandleNewMessage(c *gin.Context) {
	var data RequestData
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Process data
	respond, isNotifyNeeded := respondToQuestion(data.Question)

	// Notify admin support
	if isNotifyNeeded {
		result := advanceMessage(c, data)
		fmt.Printf(result)
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

func HandleGeneralMessage(c *gin.Context) {
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

func HandleMessageType(c *gin.Context) {
	messageType := [4]string{"sales", "partner", "support", ""}

	// Generate respondData
	respondData := map[string]interface{}{
		"items":     messageType,
		"timestamp": time.Now().UnixMilli(),
	}

	c.JSON(http.StatusCreated, respondData)
}

func advanceMessage(c *gin.Context, data RequestData) string {
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Millisecond*100)
	defer cancel()

	respch := make(chan ResultData, 2)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go notifyAdminSupport(data.Id, respch, wg)
	go sendChatGPTRequest(data.Question, respch, wg)

	wg.Wait()

	close(respch)

	for {
		select {
		case <-ctx.Done():
			return "time out"
		case <-respch:
			return "response channel"
		}
	}
}

func notifyAdminSupport(id string, respch chan ResultData, wg *sync.WaitGroup) {
	// tell admin about it
	time.Sleep(time.Millisecond * 1000)
	respch <- ResultData{Data: "Complete", Err: nil}
	wg.Done()
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

func sendChatGPTRequest(question string, respch chan ResultData, wg *sync.WaitGroup) {
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
			respch <- ResultData{Data: "error", Err: fmt.Errorf("error")}
			wg.Done()
			return
		}

		respch <- ResultData{Data: resp.Choices[0].Message.Content, Err: nil}
		wg.Done()
		return
	}

	respch <- ResultData{Data: "Api key not found", Err: fmt.Errorf("error")}
	wg.Done()
}
