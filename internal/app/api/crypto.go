package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"log"
	"math/rand"
)

func generateRandom(size int) []byte {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func generateUint32() uint32 {
	return binary.BigEndian.Uint32(generateRandom(4))
}

func getSignedUserID() string {
	return encodeID(generateUint32())
}

var secretKey = []byte("secret key")

func decodeID(msg string) (uint32, error) {
	data, err := hex.DecodeString(msg)
	if err != nil {
		log.Println("error when decoding:", err)
		return 0, err
	}
	if len(data) < 10 {
		log.Println("signature too short:", data)
		return 0, errors.New("signature too short")
	}
	id := binary.BigEndian.Uint32(data[:4])
	h := hmac.New(sha256.New, secretKey)
	h.Write(data[:4])
	sign := h.Sum(nil)

	if hmac.Equal(sign, data[4:]) {
		return id, nil
	} else {
		return 0, errors.New("signature mismatch")
	}
}

func encodeID(s uint32) string {
	h := hmac.New(sha256.New, secretKey)
	a := make([]byte, 4)
	binary.BigEndian.PutUint32(a, s)
	h.Write(a)
	sign := h.Sum(nil)
	value := append(a, sign...)
	return hex.EncodeToString(value)
}
