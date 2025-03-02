package common

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/alkurbatov/mbroker/internal/domain"
	"github.com/alkurbatov/mbroker/internal/usecase"
)

type Error struct {
	Message string `json:"message"`
}

// WriteErr записывает текст ошибки в тело ответа и выставляет HTTP код ответа.
func WriteErr(c *gin.Context, code int, err error) {
	c.Error(err) //nolint: errcheck // gin возвращает копию ошибки, добавленной в c.Errors

	c.JSON(code, Error{err.Error()})
}

func HandleError(c *gin.Context, err error) {
	if errors.Is(err, usecase.ErrNoQueue) || errors.Is(err, domain.ErrDuplicateConsumer) {
		WriteErr(c, http.StatusBadRequest, err)
		return
	}

	if errors.Is(err, domain.ErrQueueOverflow) || errors.Is(err, domain.ErrTooManyConsumers) {
		WriteErr(c, http.StatusInsufficientStorage, err)
		return
	}

	WriteErr(c, http.StatusInternalServerError, err)
}
