package utils

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func GenerateToken() string {
	id, err := gonanoid.New(21)
	if err != nil {
		panic(err)
	}
	return id
}

func GetTokenExpiration() time.Time {
	return time.Now().Add(7 * 24 * time.Hour)
}