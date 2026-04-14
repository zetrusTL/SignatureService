package device

import (
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/zetrus/signature-service/internal/domain"
	"github.com/zetrus/signature-service/internal/signing"
)

type Device struct {
	ID                 string
	Label              string
	Algorithm          domain.Algorithm
	SignatureCounter   uint64
	lastSignatureB64   string
	Signer             signing.Signer
	mu                 sync.Mutex
}

func New(id, label string, alg domain.Algorithm, s signing.Signer) *Device {
	return &Device{
		ID:        id,
		Label:     label,
		Algorithm: alg,
		Signer:    s,
	}
}

func (d *Device) Sign(data string) (signatureB64, signedData string, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var third string
	if d.SignatureCounter == 0 {
		third = base64.StdEncoding.EncodeToString([]byte(d.ID))
	} else {
		third = d.lastSignatureB64
	}

	signedData = fmt.Sprintf("%d_%s_%s", d.SignatureCounter, data, third)

	sig, err := d.Signer.SignMessage([]byte(signedData))
	if err != nil {
		return "", "", err
	}
	signatureB64 = base64.StdEncoding.EncodeToString(sig)
	d.lastSignatureB64 = signatureB64
	d.SignatureCounter++
	return signatureB64, signedData, nil
}

type Public struct {
	ID               string
	Label            string
	Algorithm        domain.Algorithm
	SignatureCounter uint64
	PublicKeyPEM     string
}

func (d *Device) Public() Public {
	return Public{
		ID:               d.ID,
		Label:            d.Label,
		Algorithm:        d.Algorithm,
		SignatureCounter: d.SignatureCounter,
		PublicKeyPEM:     d.Signer.PublicKeyPEM(),
	}
}
