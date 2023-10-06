package handler

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWT struct {
	privateKey *[]byte
	publicKey  *[]byte
}

type CustomClaims struct {
	UserId int32 `json:"user_id"`
	jwt.StandardClaims
}

func NewJWT(privateKey []byte, publicKey []byte) JWT {
	return JWT{
		privateKey: &privateKey,
		publicKey:  &publicKey,
	}
}

func (j JWT) Create(ttl time.Duration, content UserData) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(*j.privateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	var claims CustomClaims
	claims.UserId = content.Id             // Our custom data.
	claims.ExpiresAt = now.Add(ttl).Unix() // The expiration time after which the token must be disregarded.
	claims.IssuedAt = now.Unix()           // The time at which the token was issued.
	claims.NotBefore = now.Unix()          // The time before which the token must be disregarded.

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func (j JWT) Validate(token string) (*CustomClaims, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(*j.publicKey)
	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}

	tok, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := tok.Claims.(*CustomClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("validate: invalid")
	}

	return claims, nil
}
