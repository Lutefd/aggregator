package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("authorization header not found")
	}
	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("authorization header is not in the correct format")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("authorization header is not of type ApiKey")
	}
	return vals[1], nil
}
