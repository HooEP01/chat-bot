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

	// set up database
	db.StartDB()

	router := gin.Default()

	// set up templ
	router.Static("public/", "./public")
	router.LoadHTMLGlob("views/home/*")
	ginHtmlRenderer := router.HTMLRender
	router.HTMLRender = &gintemplrenderer.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}

	// api
	router.POST("/message", handle.HandleNewMessage)
	router.GET("/message/type", handle.HandleMessageType)
	router.GET("/message", handle.HandleGeneralMessage)

	// view
	// router.GET("/login", messageHandler.HandleLogin)
	router.GET("/home", handle.HandleHome)

	router.POST("/faq", handle.HandleInsertFaq)
	router.GET("/faq", handle.HandleFaqs)
	router.GET("/faq/:id", handle.HandleFaq)
	router.PUT("/faq/:id", handle.HandleUpdateFaq)
	router.DELETE("/faq/:id", handle.HandleDeleteFaq)

	router.POST("/faq-type", handle.HandleInsertFaqtype)
	router.GET("/faq-type", handle.HandleFaqtypes)
	router.GET("/faq-type/:id", handle.HandleFaqtype)
	router.PUT("/faq-type/:id", handle.HandleUpdateFaqtype)
	router.DELETE("/faq-type/:id", handle.HandleDeleteFaqtype)

	router.Run(":8080")
}
