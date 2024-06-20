package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

type metrics struct {
	logger *slog.Logger
	fileserverHits int
}

func (m *metrics) Inc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		m.fileserverHits++
	}
}

func (m *metrics) Get() int {
	return m.fileserverHits
}

func (m *metrics) Reset() {
	m.fileserverHits = 0
}