package middleware

import (
	"errors"
	"log/slog"
	"server_course/api/common"
	"server_course/db"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware(l *slog.Logger, userStore *db.DB) gin.HandlerFunc {
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

		var userID int
		var err error
		if len(strings.Split(parts[1], ".")) == 3 {
			logger.Debug("token is jwt")
			userID, err = common.ValidJWT(parts[1])
		} else {
			logger.Debug("token is refresh token")
			// var found bool
			// userID, found, err = common.ValidRefreshToken(userStore, parts[1])
			// if !found {
			// 	err = errors.New("user not found")
			// }
			err = errors.New("not supported")
		}

		if err != nil {
			logger.Debug("token not valid", slog.String("err", err.Error()))
			c.AbortWithError(401, err)
			return
		}

		logger.Debug("token valid", slog.Int("userID", userID))

		c.Set("userID", userID)
	}
}
