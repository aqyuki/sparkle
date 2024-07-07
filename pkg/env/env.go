package env

import (
	"errors"
	"os"
)

var (
	ErrEnvNotFound = errors.New("environment variable was not found")
)

func GetOrErr(key string) (string, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return "", ErrEnvNotFound
	}
	return v, nil
}
