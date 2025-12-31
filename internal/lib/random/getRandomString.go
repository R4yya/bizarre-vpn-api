package random

import (
	"crypto/rand"
	"math/big"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GetRandomString(length int) (string, error) {
	if length <= 0 {
		return "", nil
	}

	result := make([]byte, length)
	max := big.NewInt(int64(len(letters)))

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		result[i] = letters[n.Int64()]
	}

	return string(result), nil
}
