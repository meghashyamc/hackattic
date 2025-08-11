package auth

import (
	"errors"
	"log/slog"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

var tokenPattern = regexp.MustCompile(`^[0-9a-f]{16}$`)

func GetAccessToken() (string, error) {
	godotenv.Load("../../.env")
	accessToken := os.Getenv("ACCESS_TOKEN")
	if err := validateAccessToken(accessToken); err != nil {
		return "", err
	}

	return accessToken, nil
}
func validateAccessToken(token string) error {
	if len(token) == 0 {
		slog.Error("access token cannot be empty")
		return errors.New("access token cannot be empty")
	}

	if !tokenPattern.MatchString(token) {
		slog.Error("invalid token")
		return errors.New("invalid token")
	}
	return nil
}
