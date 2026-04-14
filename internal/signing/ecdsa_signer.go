package signing

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
)

type ecdsaSigner struct {
	priv *ecdsa.PrivateKey
}

func newECDSASigner() (Signer, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &ecdsaSigner{priv: priv}, nil
}

func (s *ecdsaSigner) PublicKeyPEM() string {
	b, _ := x509.MarshalPKIXPublicKey(&s.priv.PublicKey)
	block := &pem.Block{Type: "PUBLIC KEY", Bytes: b}
	return string(pem.EncodeToMemory(block))
}

func (s *ecdsaSigner) SignMessage(message []byte) ([]byte, error) {
	return signECDSA(s.priv, message)
}
