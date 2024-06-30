package common

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"server_course/db"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID, expires int) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")

	now := time.Now().UTC()
	expiresInSeconds := now.Add(time.Second * time.Duration(expires))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "chirpy",
		"iat": now.Unix(),
		"exp": expiresInSeconds.Unix(),
		"sub": strconv.Itoa(userID),
	})

	return token.SignedString([]byte(jwtSecret))
}

func ValidJWT(tokenString string) (int, error) {
	jwtSecret := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("token is not valid")
	}

	subString, err := token.Claims.GetSubject()
	if err != nil {
		return 0, err
	}

	sub, err := strconv.Atoi(subString)
	if err != nil {
		return 0, err
	}

	return sub, nil
}

func GetRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func GetAuthorizationFromHeader(c *gin.Context) (string, error) {
	token := c.GetHeader("Authorization")
	if token == "" {
		return "", errors.New("'Authorization' header not set")
	}

	parts := strings.Fields(token)
	if len(parts) != 2 {
		return "", errors.New("bad 'Authorization' header format")
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("'Authorization' header not a 'Bearer' token")
	}

	return parts[1], nil
}

func ValidRefreshToken(userStore *db.DB, refreshToken string) (int, bool, error) {
	users, err := userStore.GetUsers()
	if err != nil {
		return 0, false, err
	}

	found := false
	var userID int
	for _, user := range users {
		if strings.EqualFold(user.RefreshToken, refreshToken) {
			userID = user.ID
			found = true
			break
		}
	}

	if !found {
		return 0, false, nil
	}

	return userID, true, nil
}
