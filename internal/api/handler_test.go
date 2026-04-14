package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/zetrus/signature-service/internal/repository"
	"github.com/zetrus/signature-service/internal/service"
)

func TestAPI_CreateListGetSign(t *testing.T) {
	repo := repository.NewMemory()
	svc := service.New(repo)
	h := NewHandler(svc)
	mux := http.NewServeMux()
	h.Register(mux)

	id := uuid.MustParse("11111111-1111-1111-1111-111111111111").String()
	body := map[string]any{
		"id":        id,
		"algorithm": "RSA",
		"label":     "L",
	}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/v1/signature-devices", bytes.NewReader(b))
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusCreated {
		t.Fatalf("create: %d %s", w.Code, w.Body.String())
	}

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/v1/signature-devices", nil)
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("list: %d", w.Code)
	}

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/v1/signature-devices/"+id, nil)
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("get: %d", w.Code)
	}

	signBody := map[string]string{"data": "payload"}
	sb, _ := json.Marshal(signBody)
	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPost, "/v1/signature-devices/"+id+"/signatures", bytes.NewReader(sb))
	mux.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("sign: %d %s", w.Code, w.Body.String())
	}
	var out signResponse
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out.Signature == "" || out.SignedData == "" {
		t.Fatal("empty response fields")
	}
}
