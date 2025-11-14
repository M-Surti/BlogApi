package middleware

import (
	"fmt"
	"os"
	"strconv"

	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)



func GenerateJWT(UserID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": UserID,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func verifyJWT(tokenStirng string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenStirng, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authheader := r.Header.Get("Authorization")
		if authheader == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
		}

		tokenString := strings.TrimPrefix(authheader, "Bearer ")
		_, err := verifyJWT(tokenString)
		if err != nil {
			http.Error(w, "invalid token", http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetUserIDFromToken(r *http.Request) (uint, error) {
	authheader := r.Header.Get("Authorization")
		if authheader == "" {
			
			return 0, fmt.Errorf("authorization header missing")
		}

		tokenString := strings.TrimPrefix(authheader, "Bearer ")
		token , err := verifyJWT(tokenString)
		if err != nil {
			
			return 0, fmt.Errorf("invalid token")
		}	

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok{
		return 0, fmt.Errorf("could not pass claims")
	}

	userID, ok := claims["userID"].(string)
	if !ok{
		return 0, fmt.Errorf("userid not found in token")
	}
	userIDInt,_ := strconv.Atoi(userID)
	return uint(userIDInt), nil

}
