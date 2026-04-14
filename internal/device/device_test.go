package device

import (
	"encoding/base64"
	"sync"
	"testing"

	"github.com/google/uuid"

	"github.com/zetrus/signature-service/internal/domain"
)

type stubSigner struct {
	out []byte
}

func (s *stubSigner) SignMessage([]byte) ([]byte, error) {
	return s.out, nil
}

func (s *stubSigner) PublicKeyPEM() string {
	return ""
}

func TestSign_FirstUsesDeviceIDAsThirdPart(t *testing.T) {
	id := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8").String()
	d := New(id, "", domain.AlgorithmRSA, &stubSigner{out: []byte{1, 2, 3}})
	sigB64, signed, err := d.Sign("hello")
	if err != nil {
		t.Fatal(err)
	}
	wantThird := base64.StdEncoding.EncodeToString([]byte(id))
	wantSigned := "0_hello_" + wantThird
	if signed != wantSigned {
		t.Fatalf("signed_data: got %q want %q", signed, wantSigned)
	}
	if sigB64 != base64.StdEncoding.EncodeToString([]byte{1, 2, 3}) {
		t.Fatalf("signature: got %q", sigB64)
	}
	if d.SignatureCounter != 1 {
		t.Fatalf("counter: got %d", d.SignatureCounter)
	}
}

func TestSign_ChainsPreviousSignature(t *testing.T) {
	id := uuid.MustParse("6ba7b811-9dad-11d1-80b4-00c04fd430c8").String()
	d := New(id, "", domain.AlgorithmRSA, &stubSigner{out: []byte{9}})
	_, _, err := d.Sign("a")
	if err != nil {
		t.Fatal(err)
	}
	_, signed2, err := d.Sign("b")
	if err != nil {
		t.Fatal(err)
	}
	prev := base64.StdEncoding.EncodeToString([]byte{9})
	want := "1_b_" + prev
	if signed2 != want {
		t.Fatalf("signed_data: got %q want %q", signed2, want)
	}
	if d.SignatureCounter != 2 {
		t.Fatalf("counter: got %d", d.SignatureCounter)
	}
}

func TestSign_ConcurrentMonotonicCounter(t *testing.T) {
	id := uuid.MustParse("6ba7b812-9dad-11d1-80b4-00c04fd430c8").String()
	d := New(id, "", domain.AlgorithmRSA, &stubSigner{out: []byte{1}})
	const n = 100
	var wg sync.WaitGroup
	errs := make(chan error, n)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			_, _, err := d.Sign("x")
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	if d.SignatureCounter != n {
		t.Fatalf("counter: got %d want %d", d.SignatureCounter, n)
	}
}
