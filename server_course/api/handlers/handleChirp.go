package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"server_course/db"
	"server_course/entities"
	"strconv"
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

		c.JSON(http.StatusOK, chirps)
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

		chrip, err = chirpStore.StoreChirp(chrip)
		if err != nil {
			logger.Error("failed to StoreChirp", slog.String("err", err.Error()))
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusCreated, chrip)
	}
}
