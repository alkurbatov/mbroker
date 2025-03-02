package common

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// ParseRequest выполняет десериализацию полученного запроса в формате JSON.
func ParseRequest(c *gin.Context, v any) error {
	if err := c.ShouldBindJSON(v); err != nil {
		return fmt.Errorf("parse request body: %w", err)
	}

	return nil
}
