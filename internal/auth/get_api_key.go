package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	header := headers.Get("Authorization")
	if header == "" {
		return "", errors.New("header not found: authorization")
	}
	if !strings.HasPrefix(header, "ApiKey ") {
		return "", errors.New("api key not found")
	}
	return strings.TrimPrefix(header, "ApiKey "), nil
}
