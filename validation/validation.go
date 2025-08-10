package validation

import (
	"errors"
	"regexp"
)

var tokenPattern = regexp.MustCompile(`^[0-9a-f]{16}$`)

func ValidateAccessToken(token string) error {
	if len(token) == 0 {
		return errors.New("access token cannot be empty")
	}

	if !tokenPattern.MatchString(token) {
		return errors.New("invalid token")
	}
	return nil
}
