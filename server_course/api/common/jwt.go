package common

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(userID, expires int) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")

	expiresInSeconds := time.Now().UTC().Add(time.Duration(expires))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "chirpy",
		"iat": time.Now().UTC().Unix(),
		"exp": expiresInSeconds.Unix(),
		"sub": userID,
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

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("could not parse claims")
	}

	err = claims.Valid()
	if err != nil {
		return 0, err
	}

	issInter, ok := claims["iss"]
	if !ok {
		return 0, errors.New("issuer is wrong")
	}
	if iss, ok := issInter.(string); !ok || !strings.EqualFold(iss, "chirpy") {
		return 0, errors.New("issuer is wrong")
	}

	subInter, ok := claims["sub"]
	if !ok {
		return 0, errors.New("no subject")
	}
	sub, ok := subInter.(float64)
	if !ok {
		return 0, errors.New("no subject")
	}
	return int(sub), nil
}
