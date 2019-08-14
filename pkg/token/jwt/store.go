package jwt

import (
	"encoding/base64"
	"time"

	"elbix.dev/engine/pkg/token"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type jwtProvider struct {
	privateKey interface{}
	publicKey  interface{}
}

func loadKeyFile(data string, pub bool) (interface{}, error) {
	keyFile, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	var res interface{}
	if pub {
		res, err = jwt.ParseRSAPublicKeyFromPEM(keyFile)
	} else {
		res, err = jwt.ParseRSAPrivateKeyFromPEM(keyFile)
	}

	if err != nil {
		return nil, err
	}

	return res, nil
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
		claim[i] = data [i]
	}
	claim["exp"] = time.Now().Add(exp).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(jw.privateKey)
}

// Fetch is used to handle the verification
func (jw *jwtProvider) Fetch(token string) (map[string]interface{}, error) {
	tok, err := jwt.Parse(token, func(token *jwt.Token) (i interface{}, e error) {
		return jw.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	clm := tok.Claims.(jwt.MapClaims)
	return clm, nil
}

// NewJWTTokenProvider return a new JWT provider
func NewJWTTokenProvider(private, public string) (token.Provider, error) {
	privateKey, err := loadKeyFile(private, false)
	if err != nil {
		return nil, err
	}
	publicKey, err := loadKeyFile(public, true)
	if err != nil {
		return nil, err
	}

	return &jwtProvider{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}
