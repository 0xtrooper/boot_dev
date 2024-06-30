package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"server_course/api/handlers"
	"server_course/api/middleware"
	"server_course/db"

	"github.com/gin-gonic/gin"
)

func addRoutes(r *gin.Engine, l *slog.Logger, m middleware.Middleware, db *db.DB) {
	app := r.Group("/app")
	app.Use(m.Metrics.Inc())
	app.Static("/", "./public")

	api := r.Group("/api")
	api.GET("/healthz", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte("OK"))
	})
	api.GET("/reset", func(c *gin.Context) {
		m.Metrics.Reset()
		c.Status(http.StatusOK)
	})
	api.POST("/validate_chirp", handlers.PostValidateChirp(l))
	api.GET("/chirps", handlers.GetChirp(l, db))
	api.GET("/chirps/:chirpID", handlers.GetChirpByID(l, db))
	api.POST("/chirps", middleware.JWTMiddleware(l, db), handlers.PostChirp(l, db))
	api.DELETE("/chirps/:chirpID", middleware.JWTMiddleware(l, db), handlers.DeleteChirp(l, db))

	api.GET("/users", handlers.GetUser(l, db))
	api.GET("/users/:userID", handlers.GetUserByID(l, db))
	api.POST("/users", handlers.PostUser(l, db))
	api.PUT("/users", middleware.JWTMiddleware(l, db), handlers.PutUser(l, db))
	api.POST("/login", handlers.PostUserLogin(l, db))
	api.POST("/refresh", handlers.PostRefresh(l, db))
	api.POST("/revoke", handlers.PostRevoke(l, db))

	api.POST("/polka/webhooks", handlers.PostWebhook(l, db))

	admin := r.Group("/admin")
	admin.GET("/metrics", func(c *gin.Context) {
		responseText := fmt.Sprintf("<html>\n\n<body>\n\t<h1>Welcome, Chirpy Admin</h1>\n\t<p>Chirpy has been visited %d times!</p>\n</body>\n\n</html>", m.Metrics.Get())
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(responseText))
	})
}
