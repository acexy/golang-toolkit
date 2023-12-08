package rsa

import (
	"crypto/rand"
	"crypto/rsa"
)

func Generate(len int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, len)
	if err != nil {
		return nil, nil
	}
	publicKey := &privateKey.PublicKey
	return privateKey, publicKey
}
