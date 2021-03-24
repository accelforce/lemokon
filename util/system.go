package util

import (
	"fmt"
	"os"
)

func RequireEnv(key string) string {
	e, err := CheckEnv(key)
	if err != nil {
		panic(err)
	}
	return e
}

func CheckEnv(key string) (string, error) {
	e := os.Getenv(key)
	if len(e) == 0 {
		return "", fmt.Errorf("environment %s is not set", key)
	}
	return e, nil
}

func BoolEnv(key string) bool {
	e := os.Getenv(key)
	return e == "1" ||
		e == "t" ||
		e == "T" ||
		e == "true" ||
		e == "TRUE"
}
