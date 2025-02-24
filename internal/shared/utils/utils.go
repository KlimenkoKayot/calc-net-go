package utils

import "crypto/sha512"

func ExpressionToSHA512(expression string) [64]byte {
	hash := sha512.Sum512([]byte(expression))
	return hash
}
