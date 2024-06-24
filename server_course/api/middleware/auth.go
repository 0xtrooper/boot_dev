package middleware

import (
	"errors"
	"log/slog"
	"server_course/api/common"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware(l *slog.Logger) gin.HandlerFunc {
	logger := l.With("middleware", "JWTMiddleware")

	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		logger.Debug("read token", slog.String("token", token))
		if token == "" {
			logger.Debug("token not set")
			c.AbortWithError(401, errors.New("no jwt"))
		}

		parts := strings.Fields(token)
		if !strings.EqualFold(parts[0], "Bearer") {
			logger.Debug("token is bad", slog.String("err", "wrong bearer"))
			c.AbortWithError(401, errors.New("wrong bearer"))
		}

		userID, err := common.ValidJWT(parts[1])
		if err != nil {
			logger.Debug("token not valid", slog.String("err", err.Error()))
			c.AbortWithError(401, err)
		}

		logger.Debug("token valid", slog.Int("userID", userID))

		c.Set("userID", userID)
	}
}
