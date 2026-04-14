package signing

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

type rsaSigner struct {
	priv *rsa.PrivateKey
}

func newRSASigner() (Signer, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return &rsaSigner{priv: priv}, nil
}

func (s *rsaSigner) PublicKeyPEM() string {
	b, _ := x509.MarshalPKIXPublicKey(&s.priv.PublicKey)
	block := &pem.Block{Type: "PUBLIC KEY", Bytes: b}
	return string(pem.EncodeToMemory(block))
}

func (s *rsaSigner) SignMessage(message []byte) ([]byte, error) {
	return signRSA(s.priv, message)
}
