package user

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
    repo *Repository
}



func NewService(repo *Repository) *Service {
    return &Service{
        repo: repo,
    }
}

func (s *Service) Authenticate(ctx context.Context, username, password string) (string, error) {
    var jwtKey = []byte(os.Getenv("jwtKey"))

    user, err := s.repo.FindByUsername(ctx, username)
    if err != nil {
        return "", err
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return "", fmt.Errorf("invalid credentials")
    }

    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Subject:   user.ID.Hex(),
	}
	

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func (s *Service) Register(ctx context.Context, username, password string) error {

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    user := &User{
        Username: username,
        Password: string(hashedPassword),
    }

    return s.repo.InsertUser(ctx, user)
}
