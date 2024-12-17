package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
