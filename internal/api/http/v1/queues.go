package v1

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/alkurbatov/mbroker/internal/api/http/common"
	"github.com/alkurbatov/mbroker/internal/domain"
	"github.com/alkurbatov/mbroker/internal/usecase"
)

type subcribeRequest struct {
	URI string `binding:"required" json:"uri"`
}

type QueuesHandler struct {
	bus usecase.Bus
}

func NewQueuesHandler(bus usecase.Bus) QueuesHandler {
	return QueuesHandler{
		bus: bus,
	}
}

// Publish размещает сообщение в очереди.
func (h QueuesHandler) Publish(c *gin.Context) {
	queueName := c.Param("name")

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		common.WriteErr(c, http.StatusBadRequest, fmt.Errorf("read message: %w", err))
		return
	}

	if err = h.bus.PostMessage(queueName, domain.Message(jsonData)); err != nil {
		common.HandleError(c, err)
		return
	}
}

// Subscribe подписывает клиента на сообщения очереди.
func (h QueuesHandler) Subscribe(c *gin.Context) {
	queueName := c.Param("name")

	req := subcribeRequest{}
	if err := common.ParseRequest(c, &req); err != nil {
		common.WriteErr(c, http.StatusBadRequest, err)
		return
	}

	if err := h.bus.Subscribe(queueName, req.URI); err != nil {
		common.HandleError(c, err)
		return
	}
}
