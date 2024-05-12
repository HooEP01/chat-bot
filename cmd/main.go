package main

import (
	"time"

	"github.com/HooEP01/chat-bot/db"
	"github.com/HooEP01/chat-bot/handle"
	"github.com/a-h/templ/examples/integration-gin/gintemplrenderer"
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

	db.StartDB()

	router := gin.Default()
	router.LoadHTMLGlob("views/home/*")
	ginHtmlRenderer := router.HTMLRender
	router.HTMLRender = &gintemplrenderer.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}
	router.SetTrustedProxies(nil)

	messageHandler := handle.MessageHandler{}

	// api
	router.POST("/message", messageHandler.HandleNewMessage)
	router.GET("/message/type", messageHandler.HandleMessageType)
	router.GET("/message", messageHandler.HandleGeneralMessage)

	// view
	router.GET("/home", messageHandler.HandleHome)

	router.POST("/faq", messageHandler.HandleInsertFaq)
	router.GET("/faq", messageHandler.HandleFaqs)
	router.GET("/faq/:id", messageHandler.HandleFaq)
	router.PUT("/faq/:id", messageHandler.HandleUpdateFaq)
	router.DELETE("/faq/:id", messageHandler.HandleDeleteFaq)

	router.POST("/faq-type", messageHandler.HandleInsertFaqtype)
	router.GET("/faq-type", messageHandler.HandleFaqtypes)
	router.GET("/faq-type/:id", messageHandler.HandleFaqtype)
	router.PUT("/faq-type/:id", messageHandler.HandleUpdateFaqtype)
	router.DELETE("/faq-type/:id", messageHandler.HandleDeleteFaqtype)

	router.Run(":8080")
}
