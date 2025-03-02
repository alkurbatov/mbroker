package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/alkurbatov/mbroker/internal/usecase"
)

func Inject(router *gin.Engine, producer usecase.Producer) {
	h := NewQueuesHandler(producer)

	queues := router.Group("/v1/queues/:name")
	queues.POST("/messages", h.Publish)
}
