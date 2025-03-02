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

type QueuesHandler struct {
	producer usecase.Producer
}

func NewQueuesHandler(producer usecase.Producer) QueuesHandler {
	return QueuesHandler{
		producer: producer,
	}
}

// Publish размещает сообщение в очереди.
func (h QueuesHandler) Publish(c *gin.Context) {
	dst := c.Param("name")

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		common.WriteErr(c, http.StatusBadRequest, fmt.Errorf("read message: %w", err))
		return
	}

	if err = h.producer.PostMessage(dst, domain.Message(jsonData)); err != nil {
		common.HandleError(c, err)
		return
	}
}
