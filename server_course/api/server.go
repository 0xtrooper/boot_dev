package api

import (
	"log/slog"
	"server_course/api/middleware"
	"server_course/db"

	"github.com/gin-gonic/gin"
)

func NewServer(l *slog.Logger, m middleware.Middleware, db *db.DB) *gin.Engine {
	router := gin.Default()
	addRoutes(
		router,
		l,
		m,
		db,
	)

	return router
}
