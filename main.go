package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/sashabaranov/go-openai"
)

type RequestData struct {
	Id       string `json:"id"`
	Question string `json:"question"`
}

type QuestionAnswer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type ChatMessages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Id                string        `json:"id"`
	Object            *string       `json:"object"`
	Created           int           `json:"created"`
	Model             string        `json:"model"`
	SystemFingerprint *string       `json:"system_fingerprint"`
	Choices           *[]ChatChoice `json:"choices"`
	Usage             *ChatUsage    `json:"usage"`
}

type ChatUsage struct {
	PromptTokens     *int `json:"prompt_tokens"`
	CompletionTokens *int `json:"completion_tokens"`
	TotalTokens      *int `json:"total_tokens"`
}

type ChatChoice struct {
	Index        int           `json:"index"`
	Message      *ChatMessages `json:"message"`
	Logprobs     *string       `json:"logprobs"`
	FinishReason *string       `json:"finish_reason"`
}

var loc, _ = time.LoadLocation("Asia/Kuala_Lumpur")

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		panic("(E000) Failed to get env file.")
	}

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
			panic(err)
		}

		return resp.Choices[0].Message.Content

		// data := map[string]interface{}{
		// 	"model": "gpt-3.5-turbo",
		// 	"messages": []ChatMessages{
		// 		// {Role: "system", Content: "You are a helpful assistant."},
		// 		{Role: "user", Content: question},
		// 	},
		// }

		// body, err := json.Marshal(data)
		// if err != nil {
		// 	panic("(E001) Failed to json.Marshal(data).")
		// }

		// req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
		// if err != nil {
		// 	panic("(E002) Failed to wrap http request.")
		// }

		// req.Header.Set("Authorization", "Bearer "+apiKey)
		// req.Header.Set("Content-Type", "application/json")

		// client := &http.Client{}
		// resp, err := client.Do(req)
		// if err != nil {
		// 	panic("(E003) Failed to send request.")
		// }

		// defer resp.Body.Close()
		// if resp.StatusCode > http.StatusBadRequest {
		// 	panic("(E004) Status Code " + resp.Status)
		// }

		// var responseData ChatResponse
		// if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		// 	panic("(E005) Failed to decode response data.")
		// }

		// return responseData.Choices[0].Message.Content
	}

	return "Api key not found"
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
