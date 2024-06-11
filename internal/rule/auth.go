package rule

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)


type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Define a custom type for the context key
type contextKey string
const claimsKey contextKey = "claims"

func JWTAuth(next http.Handler) http.Handler {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var jwtKey = []byte(os.Getenv("jwtKey"))

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims := &Claims{}

        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })
        if err != nil || !token.Valid {
			println(token.Valid)
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // Add claims to context
        ctx := context.WithValue(r.Context(), claimsKey , claims)
        r = r.WithContext(ctx)

        next.ServeHTTP(w, r)
    })
}