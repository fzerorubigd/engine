package sec

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"github.com/pkg/errors"
)

// ExtractPublicFromPrivate is a useful function for dependency injection
func ExtractPublicFromPrivate(in *rsa.PrivateKey) crypto.PublicKey {
	return in.Public()
}

// ParseRSAPrivateKeyFromPEM read RSA from PEM encoded data
func ParseRSAPrivateKeyFromPEM(key []byte) (*rsa.PrivateKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, errors.New("key must be PEM encoded")
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	}

	var pkey *rsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		return nil, errors.New("not a RSA private key")
	}

	return pkey, nil
}

// ParseRSAPrivateKeyFromBase64PEM load private key from base64 encoded string
func ParseRSAPrivateKeyFromBase64PEM(private string) (*rsa.PrivateKey, error) {
	b, err := base64.StdEncoding.DecodeString(private)
	if err != nil {
		return nil, err
	}

	return ParseRSAPrivateKeyFromPEM(b)
}
