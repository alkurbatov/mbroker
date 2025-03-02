package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/alkurbatov/mbroker/internal/usecase"
)

func Inject(router *gin.Engine, bus usecase.Bus) {
	h := NewQueuesHandler(bus)

	queues := router.Group("/v1/queues/:name")
	queues.POST("/messages", h.Publish)
	queues.POST("/subscriptions", h.Subscribe)
}
