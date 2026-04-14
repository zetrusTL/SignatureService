package signing

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"testing"

	"github.com/zetrus/signature-service/internal/domain"
)

func TestRSA_SignVerify(t *testing.T) {
	s, err := NewSigner(domain.AlgorithmRSA)
	if err != nil {
		t.Fatal(err)
	}
	msg := []byte("secured payload")
	sig, err := s.SignMessage(msg)
	if err != nil {
		t.Fatal(err)
	}
	rs, ok := s.(*rsaSigner)
	if !ok {
		t.Fatal("expected rsa signer")
	}
	h := sha256.Sum256(msg)
	err = rsa.VerifyPKCS1v15(&rs.priv.PublicKey, crypto.SHA256, h[:], sig)
	if err != nil {
		t.Fatal(err)
	}
}

func TestECDSA_SignVerify(t *testing.T) {
	s, err := NewSigner(domain.AlgorithmECC)
	if err != nil {
		t.Fatal(err)
	}
	msg := []byte("secured payload")
	sig, err := s.SignMessage(msg)
	if err != nil {
		t.Fatal(err)
	}
	es, ok := s.(*ecdsaSigner)
	if !ok {
		t.Fatal("expected ecdsa signer")
	}
	h := sha256.Sum256(msg)
	if !ecdsa.VerifyASN1(&es.priv.PublicKey, h[:], sig) {
		t.Fatal("verify failed")
	}
}
