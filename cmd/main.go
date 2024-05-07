package main

import (
	"time"

	"github.com/HooEP01/chat-bot/handle"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var loc, _ = time.LoadLocation("Asia/Kuala_Lumpur")

type QuestionAnswer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		panic("(E000) Failed to get env file.")
	}

	router := gin.Default()

	messageHandler := handle.MessageHandler{}

	router.POST("/message", messageHandler.HandleNewMessage)
	router.GET("/message", messageHandler.HandleGeneralMessage)
	router.GET("/message/type", messageHandler.HandleMessageType)

	router.Run(":8080")
}
