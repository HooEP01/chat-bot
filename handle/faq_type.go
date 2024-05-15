package handle

import (
	"net/http"

	faqtype "github.com/HooEP01/chat-bot/views/faq_type"
	"github.com/a-h/templ/examples/integration-gin/gintemplrenderer"
	"github.com/gin-gonic/gin"
)

func HandleInsertFaqtype(c *gin.Context) {}
func HandleFaqtypes(c *gin.Context) {
	// messageType := [4]string{"sales", "partner", "support", ""}

	// Generate respondData
	// respondData := map[string]interface{}{
	// 	"items":     messageType,
	// 	"timestamp": time.Now().UnixMilli(),
	// }

	r := gintemplrenderer.New(c.Request.Context(), http.StatusOK, faqtype.Hello("This is text"))
	c.Render(http.StatusOK, r)
}
func HandleFaqtype(c *gin.Context)       {}
func HandleUpdateFaqtype(c *gin.Context) {}
func HandleDeleteFaqtype(c *gin.Context) {}
