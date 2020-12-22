package utils

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/FaranushKarimov/crud/pkg/types"
)

//GenerateTokenStr ...
func GenerateTokenStr() (string, error) {

	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		return "", types.ErrInternal
	}

	return hex.EncodeToString(buffer), nil
}
