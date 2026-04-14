package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/zetrus/signature-service/internal/device"
	"github.com/zetrus/signature-service/internal/domain"
	"github.com/zetrus/signature-service/internal/service"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/v1/signature-devices", h.devicesCollection)
	mux.HandleFunc("/v1/signature-devices/", h.deviceSubpaths)
}

func (h *Handler) devicesCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createDevice(w, r)
	case http.MethodGet:
		h.listDevices(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) deviceSubpaths(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(r.URL.Path, "/v1/signature-devices/")
	p = strings.Trim(p, "/")
	parts := strings.Split(p, "/")
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}
	id := parts[0]
	if len(parts) == 1 {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.getDevice(w, r, id)
		return
	}
	if len(parts) == 2 && parts[1] == "signatures" {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.sign(w, r, id)
		return
	}
	http.NotFound(w, r)
}

type createDeviceRequest struct {
	ID        string `json:"id"`
	Algorithm string `json:"algorithm"`
	Label     string `json:"label"`
}

type deviceResponse struct {
	ID               string `json:"id"`
	Label            string `json:"label"`
	Algorithm        string `json:"algorithm"`
	SignatureCounter uint64 `json:"signature_counter"`
	PublicKeyPEM     string `json:"public_key_pem"`
}

type signRequest struct {
	Data string `json:"data"`
}

type signResponse struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *Handler) createDevice(w http.ResponseWriter, r *http.Request) {
	var req createDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}
	alg, err := domain.ParseAlgorithm(req.Algorithm)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}
	pub, err := h.svc.CreateDevice(r.Context(), req.ID, req.Label, alg)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toDeviceResponse(pub))
}

func (h *Handler) listDevices(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.ListDevices(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
		return
	}
	out := make([]deviceResponse, 0, len(list))
	for _, d := range list {
		out = append(out, toDeviceResponse(d))
	}
	writeJSON(w, http.StatusOK, out)
}

func (h *Handler) getDevice(w http.ResponseWriter, r *http.Request, id string) {
	pub, err := h.svc.GetDevice(r.Context(), id)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toDeviceResponse(pub))
}

func (h *Handler) sign(w http.ResponseWriter, r *http.Request, id string) {
	var req signRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json"})
		return
	}
	sig, signed, err := h.svc.SignTransaction(r.Context(), id, req.Data)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, signResponse{Signature: sig, SignedData: signed})
}

func (h *Handler) writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "not found"})
	case errors.Is(err, domain.ErrConflict):
		writeJSON(w, http.StatusConflict, errorResponse{Error: "already exists"})
	case errors.Is(err, domain.ErrInvalidID):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	case errors.Is(err, domain.ErrUnknownAlgorithm):
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
	}
}

func toDeviceResponse(d device.Public) deviceResponse {
	return deviceResponse{
		ID:               d.ID,
		Label:            d.Label,
		Algorithm:        string(d.Algorithm),
		SignatureCounter: d.SignatureCounter,
		PublicKeyPEM:     d.PublicKeyPEM,
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
