package handle

import (
	"net/http"
	"time"

	"github.com/HooEP01/chat-bot/types"
	"github.com/gin-gonic/gin"
)

func HandleInsertFaq(c *gin.Context) {}
func HandleFaqs(c *gin.Context) {
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
func HandleFaq(c *gin.Context)       {}
func HandleUpdateFaq(c *gin.Context) {}
func HandleDeleteFaq(c *gin.Context) {}
