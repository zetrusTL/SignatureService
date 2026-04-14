package signing

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"

	"github.com/zetrus/signature-service/internal/domain"
)

type Signer interface {
	SignMessage(message []byte) ([]byte, error)
	PublicKeyPEM() string
}

func NewSigner(alg domain.Algorithm) (Signer, error) {
	switch alg {
	case domain.AlgorithmRSA:
		return newRSASigner()
	case domain.AlgorithmECC:
		return newECDSASigner()
	default:
		return nil, domain.ErrUnknownAlgorithm
	}
}

func hashMessage(message []byte) []byte {
	h := sha256.Sum256(message)
	return h[:]
}

func signRSA(priv *rsa.PrivateKey, message []byte) ([]byte, error) {
	return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashMessage(message))
}

func signECDSA(priv *ecdsa.PrivateKey, message []byte) ([]byte, error) {
	h := hashMessage(message)
	return ecdsa.SignASN1(rand.Reader, priv, h)
}
