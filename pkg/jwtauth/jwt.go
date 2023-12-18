package jwtauth

import (
	"crypto/rsa"
	"fmt"
	"math"
	"time"

	"github.com/TGRZiminiar/Clean-Architecture-Go/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type (
	AuthFactory interface {
		SignToken() (string, error)
	}

	Claims struct {
		UserId   uuid.UUID `json:"userId"`
		Email    string    `json:"email"`
		Username string    `json:"username"`
	}

	AuthMapClaims struct {
		*Claims
		jwt.RegisteredClaims
		isErr error
	}

	authConcrete struct {
		PrivateKey *rsa.PrivateKey
		Claims     *AuthMapClaims `json:"claims"`
	}
)

func NewAccessToken(cfg *config.Jwt, claims *Claims, expiredAt int64, subject string) *authConcrete {

	// privateKeyPem is a string
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.PrivateKeyPem))
	if err != nil {
		return &authConcrete{
			PrivateKey: nil,
			Claims: &AuthMapClaims{
				Claims:           nil,
				RegisteredClaims: jwt.RegisteredClaims{},
				isErr:            fmt.Errorf("reading private key errors : %s", err.Error()),
			},
		}
	}
	return &authConcrete{
		PrivateKey: privateKey,
		Claims: &AuthMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "mix.com",
				Subject:   subject,
				Audience:  []string{"mix.com"},
				ExpiresAt: jwtTimeDurationCal(expiredAt),
				NotBefore: jwt.NewNumericDate(now()),
				IssuedAt:  jwt.NewNumericDate(now()),
			},
			isErr: nil,
		},
	}
}

func (a *authConcrete) SignToken() (string, error) {
	if a.Claims.isErr != nil {
		return "", a.Claims.isErr
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, a.Claims)
	token, err := jwtToken.SignedString(a.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("signing token error: %s", err.Error())
	}
	return token, nil
}

func ParseToken(tokenString string, cfg *config.Jwt) (*AuthMapClaims, error) {
	// publicKeyPem is a string
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cfg.PublicKeyPem))
	if err != nil {
		return nil, fmt.Errorf("reading public key errors: %s", err.Error())
	}

	token, err := jwt.ParseWithClaims(tokenString, &AuthMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("parsing token error: %s", err.Error())
	}

	if claims, ok := token.Claims.(*AuthMapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func now() time.Time {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return time.Now().In(loc)
}

func jwtTimeDurationCal(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(now().Add(time.Duration(t * int64(math.Pow10(9)))))
}
