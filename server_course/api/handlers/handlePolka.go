package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"server_course/db"
	"strings"

	"github.com/gin-gonic/gin"
)

type PolkaWebhookBody struct {
	Event string `json:"event"`
	Data  struct {
		UserID int `json:"user_id"`
	} `json:"data"`
}

func PostWebhook(l *slog.Logger, userStore *db.DB) gin.HandlerFunc {
	logger := l.With("handler", "PostWebhook")
	expectedApiKey := os.Getenv("POLKA_KEY")

	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.Status(http.StatusUnauthorized)
			return
		}

		parts := strings.Fields(token)
		if len(parts) != 2 {
			c.Status(http.StatusUnauthorized)
			return
		}

		if !strings.EqualFold(parts[0], "ApiKey") {
			c.Status(http.StatusUnauthorized)
			return
		}

		if !strings.EqualFold(parts[1], expectedApiKey) {
			c.Status(http.StatusUnauthorized)
			return
		}

		logger.Debug("request has the correct api key")

		var body PolkaWebhookBody
		err := decode(c, &body)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if !strings.EqualFold(body.Event, "user.upgraded") {
			c.Status(http.StatusNoContent)
			return
		}

		_, err = userStore.UpdateUserRedStatus(body.Data.UserID, true)
		if err != nil {
			if errors.Is(err, db.ErrDoesNotExist) {
				c.Status(http.StatusNotFound)
				return
			}
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		logger.Debug("user updated", slog.Int("userID", body.Data.UserID))
		c.Status(http.StatusNoContent)
	}
}
