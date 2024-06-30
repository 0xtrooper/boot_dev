package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"server_course/db"
	"server_course/entities"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func PostValidateChirp(l *slog.Logger) gin.HandlerFunc {
	logger := l.With("handler", "PostValidateChirp")

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*250)
		defer cancel()

		var chrip entities.Chirp
		problems, err := decodeValid(ctx, c, &chrip)
		if len(problems) > 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, problems)
			return
		}
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		logger.Debug("new validate_chirp", slog.String("body", chrip.Body))

		c.JSON(http.StatusOK, gin.H{
			"cleaned_body": chrip.Body,
		})
	}
}

func GetChirp(l *slog.Logger, chirpStore *db.DB) gin.HandlerFunc {
	logger := l.With("handler", "GetChirp")

	return func(c *gin.Context) {
		chirps, err := chirpStore.GetChirpsSlice()
		if err != nil {
			logger.Error("failed to GetChirps", slog.String("err", err.Error()))
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		direction := c.Query("sort")
		if strings.EqualFold(direction, "desc") {
			n := len(chirps)
			temp := make([]entities.Chirp, n)
			for i := 0; i < n; i++ {
				temp[n-i-1] = chirps[i]
			}
			chirps = temp
		}

		requestedAuthorIDString := c.Query("author_id")
		if requestedAuthorIDString == "" {
			c.JSON(http.StatusOK, chirps)
			return
		}

		requestedAuthorID, err := strconv.Atoi(requestedAuthorIDString)
		if err != nil {
			logger.Debug("user author id is bad", slog.String("err", err.Error()))
			c.JSON(http.StatusOK, chirps)
			return
		}

		filteredChirps := []entities.Chirp{}
		for _, chirp := range chirps {
			if chirp.AuthorID == requestedAuthorID {
				filteredChirps = append(filteredChirps, chirp)
			}
		}

		c.JSON(http.StatusOK, filteredChirps)
	}
}

func GetChirpByID(l *slog.Logger, chirpStore *db.DB) gin.HandlerFunc {
	logger := l.With("handler", "GetChirp")

	return func(c *gin.Context) {
		chirpIDString := c.Param("chirpID")
		if chirpIDString == "" {
			logger.Debug("chirpID not set")
			c.AbortWithError(http.StatusBadRequest, errors.New("chirpID not set"))
			return
		}

		chirpID, err := strconv.Atoi(chirpIDString)
		if err != nil {
			logger.Debug("chirpID not an int", slog.String("err", err.Error()))
			c.AbortWithError(http.StatusBadRequest, errors.New("chirpID not an int"))
			return
		}

		chirps, err := chirpStore.GetChirp(chirpID)
		if err != nil {
			if errors.Is(err, db.ErrDoesNotExist) {
				c.AbortWithError(http.StatusNotFound, err)
			}

			logger.Error("failed to GetChirps", slog.String("err", err.Error()))
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, chirps)
	}
}

func PostChirp(l *slog.Logger, chirpStore *db.DB) gin.HandlerFunc {
	logger := l.With("handler", "PostChirp")

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*250)
		defer cancel()

		var chrip entities.Chirp
		problems, err := decodeValid(ctx, c, &chrip)
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
			c.AbortWithError(http.StatusInternalServerError, errors.New("userID is not set"))
			return
		}

		chrip.AuthorID = userID

		chrip, err = chirpStore.StoreChirp(chrip)
		if err != nil {
			logger.Error("failed to StoreChirp", slog.String("err", err.Error()))
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, chrip)
	}
}

func DeleteChirp(l *slog.Logger, chirpStore *db.DB) gin.HandlerFunc {
	logger := l.With("handler", "DeleteChirp")

	return func(c *gin.Context) {
		chirpIDString := c.Param("chirpID")
		if chirpIDString == "" {
			logger.Debug("chirpID not set")
			c.AbortWithError(http.StatusBadRequest, errors.New("chirpID not set"))
			return
		}

		chirpID, err := strconv.Atoi(chirpIDString)
		if err != nil {
			logger.Debug("chirpID not an int", slog.String("err", err.Error()))
			c.AbortWithError(http.StatusBadRequest, errors.New("chirpID not an int"))
			return
		}

		chirp, err := chirpStore.GetChirp(chirpID)
		if err != nil {
			if errors.Is(err, db.ErrDoesNotExist) {
				c.AbortWithError(http.StatusNotFound, err)
			}

			logger.Error("failed to GetChirps", slog.String("err", err.Error()))
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		userID := c.GetInt("userID")
		if userID == 0 {
			logger.Error("userID is not set")
			c.AbortWithError(http.StatusInternalServerError, errors.New("userID is not set"))
			return
		}

		if chirp.AuthorID != userID {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		err = chirpStore.DeleteChirp(chirpID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}
