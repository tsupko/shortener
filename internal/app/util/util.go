package util

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	alphabet         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	alphabetLength   = len(alphabet)
	shortURLLength   = 8
	UserIDCookieName = "userID"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateUniqueID() string {
	result := make([]byte, shortURLLength)
	for symbol := range result {
		result[symbol] = alphabet[rand.Intn(alphabetLength)]
	}
	return string(result)
}

func ReadRequestBody(r *http.Request) (string, error) {
	defer func() {
		_ = r.Body.Close()
	}()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading request body: %v\n", err)
	}
	return string(body), err
}
