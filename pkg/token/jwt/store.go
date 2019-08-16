package jwt

import (
	"crypto/rsa"
	"time"

	"elbix.dev/engine/pkg/token"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type jwtProvider struct {
	privateKey *rsa.PrivateKey
}

func (jw *jwtProvider) Delete(token string) {
	// NO OP on jwt // TODO: HOW?
}

// Store return a new JWT token for the user
func (jw *jwtProvider) Store(data map[string]interface{}, exp time.Duration) (string, error) {
	if exp <= 0 {
		return "", errors.New("negative expiration")
	}
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	claim := make(jwt.MapClaims, len(data)+1)
	for i := range data {
		claim[i] = data[i]
	}
	claim["exp"] = time.Now().Add(exp).Unix()

	tk := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)

	// Sign and get the complete encoded token as a string using the secret
	return tk.SignedString(jw.privateKey)
}

// Fetch is used to handle the verification
func (jw *jwtProvider) Fetch(token string) (map[string]interface{}, error) {
	tok, err := jwt.Parse(token, func(token *jwt.Token) (i interface{}, e error) {
		return jw.privateKey.Public(), nil
	})
	if err != nil {
		return nil, err
	}

	clm := tok.Claims.(jwt.MapClaims)
	return clm, nil
}

// NewJWTTokenProvider return a new JWT provider
func NewJWTTokenProvider(private *rsa.PrivateKey) token.Provider {
	return &jwtProvider{
		privateKey: private,
	}
}
