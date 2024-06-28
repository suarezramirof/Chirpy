package auth

import (
	"fmt"
	"net/http"
	"strings"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"errors"
	"time"
	"crypto/rand"
	"encoding/hex"
)

func HashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func CompareHashAndPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GetBearerToken(headers http.Header) (string, error) {
	authorizationHeader := headers.Get("Authorization")
	if authorizationHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}
	splitToken := strings.Split(authorizationHeader, " ")
	if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
		return "", fmt.Errorf("invalid authorization header")
	}
	return splitToken[1], nil
}

func CheckJWT(token string, jwtSecret string) (string, error) {
	tok, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(jwtTok *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	// Check if the token is valid
	if err != nil {
		return "", err
	}
	
	userId, err := tok.Claims.GetSubject()

	if err != nil {
		return "", errors.New("invalid user token")
	}
	return userId, nil
}

func MakeJWT(userId int, jwtSecret string) (string, error) {
	expiresIn := time.Hour
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  &jwt.NumericDate{Time: time.Now().UTC()},
		ExpiresAt: &jwt.NumericDate{Time: time.Now().UTC().Add(expiresIn)},
		Subject:   fmt.Sprintf("%d", userId),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tok, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tok, nil
}

func MakeRefreshToken() (string, error) {
	// Generate a random token
	token := make([]byte, 32)
	_, err := rand.Read(token);
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	hexToken := hex.EncodeToString(token)
	return hexToken, nil
}