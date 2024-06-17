package middleware

import (
    "context"
    "log"
    "os"
    "strings"
    "errors"

    "github.com/dgrijalva/jwt-go"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
    "github.com/joho/godotenv"
)

// Define a custom type for the context key
type contextKey string
const ClaimsKey contextKey = "claims"


// Claims structure
type Claims struct {
    ID string `json:"id"`
    jwt.StandardClaims
}

// AuthInterceptor is a server interceptor for handling JWT authentication
func AuthInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {

    // Skip authentication for Login and Register methods
    if info.FullMethod == "/userpb.UserService/Login" || info.FullMethod == "/userpb.UserService/Register" {
        return handler(ctx, req)
    }

    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    var jwtKey = []byte(os.Getenv("jwtKey"))

    // Extract the "authorization" header
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, errors.New("missing metadata")
    }

    authHeader, ok := md["authorization"]
    if !ok || len(authHeader) == 0 {
        return nil, errors.New("authorization token required")
    }

    tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")

    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil || !token.Valid {
        return nil, errors.New("invalid token")
    }

    // Add claims to context
    ctx = context.WithValue(ctx, ClaimsKey, claims)

    // Continue processing the request
    return handler(ctx, req)
}