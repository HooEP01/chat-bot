package handle

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h MessageHandler) HandleInsertFaqtype(c *gin.Context) {}
func (h MessageHandler) HandleFaqtypes(c *gin.Context) {
	messageType := [4]string{"sales", "partner", "support", ""}

	// Generate respondData
	respondData := map[string]interface{}{
		"items":     messageType,
		"timestamp": time.Now().UnixMilli(),
	}

	c.JSON(http.StatusCreated, respondData)
}
func (h MessageHandler) HandleFaqtype(c *gin.Context)       {}
func (h MessageHandler) HandleUpdateFaqtype(c *gin.Context) {}
func (h MessageHandler) HandleDeleteFaqtype(c *gin.Context) {}
