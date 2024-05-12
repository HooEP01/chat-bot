package handle

import (
	"net/http"

	"github.com/HooEP01/chat-bot/views/home"
	"github.com/a-h/templ/examples/integration-gin/gintemplrenderer"
	"github.com/gin-gonic/gin"
)

func (h MessageHandler) HandleHome(c *gin.Context) {
	r := gintemplrenderer.New(c.Request.Context(), http.StatusOK, home.Hello("This is text"))
	c.Render(http.StatusOK, r)
}
