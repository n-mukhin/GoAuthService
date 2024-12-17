package tokens

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID    string `json:"user_id"`
	IPAddress string `json:"ip_address"`
	jwt.StandardClaims
}

func GenerateAccessToken(secret, userID, ip string, duration time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:    userID,
		IPAddress: ip,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(duration).Unix(),
			IssuedAt:  now.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(secret))
}

func ValidateAccessToken(secret, tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}
