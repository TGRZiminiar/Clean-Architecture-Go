package jwtauth

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
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
	}

	authConcrete struct {
		PrivateKey *rsa.PrivateKey
		PublicKey  *rsa.PublicKey
		Secret     []byte
		Claims     *AuthMapClaims `json:"claims"`
	}

	accessToken  struct{ *authConcrete }
	refreshToken struct{ *authConcrete }
	apiKey       struct{ *authConcrete }
)

func (a *authConcrete) SignToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, a.Claims)
	ss, err := token.SignedString(a.Secret)
	if err != nil {
		return "", fmt.Errorf("Singtoken error: %s", err.Error())
	}
	return ss, nil
}

func now() time.Time {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	return time.Now().In(loc)
}

// Note that: t is a second unit
// base on current time
func jwtTimeDurationCal(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(now().Add(time.Duration(t * int64(math.Pow10(9)))))
}

// base on previous time in cookie
func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func NewAccessToken(expiredAt int64, claims *Claims) AuthFactory {

	privatePemString := `-----BEGIN PRIVATE KEY-----
	MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC8LGawNydfn5TQ
	lJtcQ//vq/epAAvbhfRApHJvc/Tkx7CGTpT9XLylMXBzIuONVnYowWAN5PQObjLh
	8fMzeotSQZHdYjzEekQzpWeTmqXOuk+ZwnhrWpo40MSpY8a3JOCHc88WqkBFhADJ
	mpjjQ88UlL/WLhBlpu9Bp/VbIh09aFzyxjyPoY/E/MdwP1SeTZCsd7sgPNrIolFA
	lMrkqIp8XVu+MLPSOc8dFm08N1Bpfky5Sxfj+3brdyyCDEmgVWSsyUAYdp+trUNV
	nBpcrnDFPfleFtbNl1Wbbo2oN74k2TakbYKpUPE1/Z0KaMJoYZA5F+vmaAvTtEOh
	id7E0V1PAgMBAAECggEAE0IqvaeHZ2lPs7aMh/mK/Mk0LvRgQkwICYL1wKO00U1W
	qwg+hgezDufBGHzeSR2A46mW3mqA8pEjXViD3tA7Pg8+kdJ3YaxvGny33bVwJCOR
	TgKcjHAuDy6JUhocwQBh3RBkhhn6kKLH6do2rNvIb9TMJWwlesyGTsk/z99CmtbM
	ftLNJCVhmIyXFYLjkolyQn7qezJgJbXCUjr+XTGZZ/WV6NhCDbmt4tzXGQT679X3
	NIIpEVA+Va2up2Ua4eNc+vDfhRzNgcIE2hq4DaRyc/zSY22ItnSFxv3qQF1a3d+L
	S7ajQrqpiPC8s+yrbZQowCTUGITLV9xByW1Dghh78QKBgQDrHnyhz77jiwtt5Toe
	2c93DUOB4bu2gZkZi7mjRKCVpc+QhsBDe7hTCEruJuOiFaVfU369b12ZccK7+9ze
	OqlgaW8VJBUVR7FBs4lClUe9Sb8qIH32r4jMFBygfF44eaaDl0/Ax8qIOL5WM2A6
	6f0+1vdllwvVfLbIE9MsdqfAxQKBgQDM4pbBIpDaKCZW/qR6fYoA10PPNiA42+rC
	TfvHGxvh2QJ2kh43xliIjc7bTPI6plsWrStuZh0tyDoQGSNZHYbPBV/TCx3Ib9L9
	1k18sr8RNFQZTQUJL4qFUsNJUlZVLMxQ3Nmavt9HV0kwOEGga42Tp8kWU7PmeBVq
	CL//wh1fAwKBgQCX/cdX2zJtai9jRXIDC47gSUTAq6prWvAb4YWKFA0zcFLz/QhB
	F7OaiZvWxHEXEKMtMo6V6244iZ/3YePwDT/9QWs74W13qjbeYC91SYdsyEW59/M2
	C1eFheLTpFJMc+e+3YwC9aTp1rTEiMXGkAjUHKcllzVhNxP510cGUVY0eQKBgHd2
	o1pn0jgx8vEEt1jovD/zRImc0Lr2l/LFz8nvp5lPlJ0YY+A3mcW9keDTA+Zou3IE
	dO+BQQBB4IEkdzTt/33Ub2Q59hq6ATea7kGIY9ofPe4mt4n8m3NTp6SoCsjNPzDj
	JUqSgtQxM+6WzsVAESQIUDrhgWMfn7Tc9z6kq8WLAoGBAJMxPi8sqsrJDKDpZyrw
	/Qby4r8dKcOsOjD0U5t2tFyD77iuO6JSvxTvUCBw5aIGeikAaogm9F3hnXOJ+k+J
	gEd4l7ED/qDnkvYu7E1qjwwd3KbyEy4cSJIx0csEKcYSTSyIDVUP3/0o40KD8y3i
	n1unMlXyGaDjpp+REUr3I+fE
	-----END PRIVATE KEY-----`

	block, _ := pem.Decode([]byte(privatePemString))
	parseResult, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
	privateKey := parseResult.(*rsa.PrivateKey)

	publicPemString := `-----BEGIN PUBLIC KEY-----
	MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvCxmsDcnX5+U0JSbXEP/
	76v3qQAL24X0QKRyb3P05Mewhk6U/Vy8pTFwcyLjjVZ2KMFgDeT0Dm4y4fHzM3qL
	UkGR3WI8xHpEM6Vnk5qlzrpPmcJ4a1qaONDEqWPGtyTgh3PPFqpARYQAyZqY40PP
	FJS/1i4QZabvQaf1WyIdPWhc8sY8j6GPxPzHcD9Unk2QrHe7IDzayKJRQJTK5KiK
	fF1bvjCz0jnPHRZtPDdQaX5MuUsX4/t263csggxJoFVkrMlAGHafra1DVZwaXK5w
	xT35XhbWzZdVm26NqDe+JNk2pG2CqVDxNf2dCmjCaGGQORfr5mgL07RDoYnexNFd
	TwIDAQAB
	-----END PUBLIC KEY-----`

	blockPub, _ := pem.Decode([]byte(publicPemString))
	parseResultPub, _ := x509.ParsePKCS8PrivateKey(blockPub.Bytes)
	publicKey := parseResultPub.(*rsa.PublicKey)

	return &accessToken{
		authConcrete: &authConcrete{
			PrivateKey: privateKey,
			PublicKey:  publicKey,
			Claims: &AuthMapClaims{
				Claims: claims,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "mix.com",
					Subject:   "access-token",
					Audience:  []string{"mix.com"},
					ExpiresAt: jwtTimeDurationCal(expiredAt),
					NotBefore: jwt.NewNumericDate(now()),
					IssuedAt:  jwt.NewNumericDate(now()),
				},
			},
		},
	}
}

func NewRefreshToken(secret string, expiredAt int64, claims *Claims) AuthFactory {
	return &refreshToken{
		authConcrete: &authConcrete{
			Secret: []byte(secret),
			Claims: &AuthMapClaims{
				Claims: claims,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "mix.com",
					Subject:   "refresh-token",
					Audience:  []string{"mix.com"},
					ExpiresAt: jwtTimeDurationCal(expiredAt),
					NotBefore: jwt.NewNumericDate(now()),
					IssuedAt:  jwt.NewNumericDate(now()),
				},
			},
		},
	}
}

func ReloadToken(secret string, expiredAt int64, claims *Claims) (string, error) {
	obj := &refreshToken{
		authConcrete: &authConcrete{
			Secret: []byte(secret),
			Claims: &AuthMapClaims{
				Claims: claims,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "mix.com",
					Subject:   "refresh-token",
					Audience:  []string{"mix.com"},
					ExpiresAt: jwtTimeRepeatAdapter(expiredAt),
					NotBefore: jwt.NewNumericDate(now()),
					IssuedAt:  jwt.NewNumericDate(now()),
				},
			},
		},
	}

	return obj.SignToken()
}

func NewApiKey(secret string) AuthFactory {
	return &apiKey{
		authConcrete: &authConcrete{
			Secret: []byte(secret),
			Claims: &AuthMapClaims{
				Claims: &Claims{},
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "mix.com",
					Subject:   "api-key",
					Audience:  []string{"mix.com"},
					ExpiresAt: jwtTimeDurationCal(31560000),
					NotBefore: jwt.NewNumericDate(now()),
					IssuedAt:  jwt.NewNumericDate(now()),
				},
			},
		},
	}
}

func ParseToken(secret string, tokenString string) (*AuthMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error: unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errors.New("error: token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("error: token is expired")
		} else {
			return nil, errors.New("error: token is invalid")
		}
	}

	if claims, ok := token.Claims.(*AuthMapClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("error: claims type is invalid")
	}
}

// Apikey  generator
var apiKeyInstant string

// Work only once
var once sync.Once

func SetApiKey(secret string) {
	once.Do(func() {
		apiKeyInstant, _ = NewApiKey(secret).SignToken()
	})
}

func SetApiKeyInContext(pctx *context.Context) {
	*pctx = metadata.NewOutgoingContext(*pctx, metadata.Pairs("auth", apiKeyInstant))
}
