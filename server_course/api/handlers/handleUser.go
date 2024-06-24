package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"server_course/api/common"
	"server_course/db"
	"server_course/entities"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetUser(l *slog.Logger, userStore *db.DB) gin.HandlerFunc {
	logger := l.With("handler", "GetUser")

	return func(c *gin.Context) {
		users, err := userStore.GetUsers()
		if err != nil {
			logger.Error("failed to GetUsers", slog.String("err", err.Error()))
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func GetUserByID(l *slog.Logger, userStore *db.DB) gin.HandlerFunc {
	logger := l.With("handler", "GetUser")

	return func(c *gin.Context) {
		userIDString := c.Param("userID")
		if userIDString == "" {
			logger.Debug("userID not set")
			c.AbortWithError(http.StatusBadRequest, errors.New("userID not set"))
			return
		}

		userID, err := strconv.Atoi(userIDString)
		if err != nil {
			logger.Debug("userID not an int", slog.String("err", err.Error()))
			c.AbortWithError(http.StatusBadRequest, errors.New("userID not an int"))
			return
		}

		user, err := userStore.GetUser(userID)
		if err != nil {
			if errors.Is(err, db.ErrDoesNotExist) {
				c.AbortWithError(http.StatusNotFound, err)
			}

			logger.Error("failed to GetUser", slog.String("err", err.Error()))
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func PostUser(l *slog.Logger, userStore *db.DB) gin.HandlerFunc {
	logger := l.With("handler", "PostUser")

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*250)
		defer cancel()

		var user entities.User
		problems, err := decodeValid(ctx, c, &user)
		if len(problems) > 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, problems)
			return
		}
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		user, err = userStore.StoreUser(user)
		if err != nil {
			logger.Error("failed to StoreUser", slog.String("err", err.Error()))
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		userByte, err := user.MarshalJSONCustom()
		if err != nil {
			logger.Error("failed to MarshalJSONCustom", slog.String("err", err.Error()))
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Data(http.StatusCreated, "application/json", userByte)
	}
}

func PostUserLogin(l *slog.Logger, userStore *db.DB) gin.HandlerFunc {
	logger := l.With("handler", "PostUserLogin")

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*250)
		defer cancel()

		var user entities.User
		problems, err := decodeValid(ctx, c, &user)
		if len(problems) > 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, problems)
			return
		}
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		storedUser, err := userStore.GetUserByEmail(user.Email)
		if err != nil {
			logger.Error("failed to GetUserByEmail", slog.String("user email", user.Email), slog.String("err", err.Error()))
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		validPassword, err := user.ValidPassword(storedUser.Password)
		if err != nil {
			logger.Error("failed to ValidPassword",
				slog.String("expected password", storedUser.Password),
				slog.String("got password", user.Password),
				slog.String("err", err.Error()),
			)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if !validPassword {
			logger.Debug("failed to ValidPassword")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if user.ExpiresInSeconds > 0 {
			storedUser.Token, err = common.GenerateJWT(storedUser.ID, user.ExpiresInSeconds)
		} else {
			storedUser.Token, err = common.GenerateJWT(storedUser.ID, 24*60*60) // 24h
		}

		if err != nil {
			logger.Error("failed to GenerateJWT", slog.Int("userID", storedUser.ID), slog.String("err", err.Error()))
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		userByte, err := storedUser.MarshalJSONCustom()
		if err != nil {
			logger.Error("failed to MarshalJSONCustom", slog.String("err", err.Error()))
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Data(http.StatusOK, "application/json", userByte)
	}
}

func PutUser(l *slog.Logger, userStore *db.DB) gin.HandlerFunc {
	logger := l.With("handler", "PutUser")

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*250)
		defer cancel()

		var user entities.User
		problems, err := decodeValid(ctx, c, &user)
		if len(problems) > 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, problems)
			return
		}
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		userID := c.GetInt("userID")
		if userID == 0 {
			logger.Error("userID is not set")
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		user.ID = userID

		user, err = userStore.UpdateUser(user)
		if err != nil {
			logger.Error("failed to StoreUser", slog.String("err", err.Error()))
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		userByte, err := user.MarshalJSONCustom()
		if err != nil {
			logger.Error("failed to MarshalJSONCustom", slog.String("err", err.Error()))
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Data(http.StatusOK, "application/json", userByte)
	}
}
