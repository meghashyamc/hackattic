package validation

import (
	"errors"
	"log/slog"
	"regexp"
)

var tokenPattern = regexp.MustCompile(`^[0-9a-f]{16}$`)

func ValidateAccessToken(token string) error {
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
