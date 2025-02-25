package utils

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
)

func ExpressionToSHA512(expression string) [64]byte {
	hash := sha512.Sum512([]byte(expression))
	return hash
}

func EncodeToString(hash [64]byte) string {
	return base64.URLEncoding.EncodeToString(hash[:])
}

func EncodedToSHA512(encoded string) ([64]byte, error) {
	hash, err := base64.URLEncoding.DecodeString(encoded)
	if len(hash) != 64 {
		return [64]byte{}, ErrInvalidBase64Decode
	}
	return [64]byte(hash), err
}

func ErrorResponse(err error) []byte {
	type Error struct {
		Error string `json:"error"`
	}
	data, _ := json.Marshal(Error{
		Error: err.Error(),
	})
	return data
}
