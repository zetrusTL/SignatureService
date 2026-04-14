package service

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/zetrus/signature-service/internal/device"
	"github.com/zetrus/signature-service/internal/domain"
	"github.com/zetrus/signature-service/internal/repository"
	"github.com/zetrus/signature-service/internal/signing"
)

type Service struct {
	repo repository.DeviceRepository
}

func New(repo repository.DeviceRepository) *Service {
	return &Service{repo: repo}
}

func normalizeDeviceID(id string) (string, error) {
	if strings.TrimSpace(id) == "" {
		return uuid.New().String(), nil
	}
	return strictUUID(id)
}

func strictUUID(id string) (string, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return "", domain.ErrInvalidID
	}
	return parsed.String(), nil
}

func (s *Service) CreateDevice(ctx context.Context, id, label string, alg domain.Algorithm) (device.Public, error) {
	nid, err := normalizeDeviceID(id)
	if err != nil {
		return device.Public{}, err
	}
	signer, err := signing.NewSigner(alg)
	if err != nil {
		return device.Public{}, err
	}
	d := device.New(nid, label, alg, signer)
	if err := s.repo.Create(ctx, d); err != nil {
		return device.Public{}, err
	}
	return d.Public(), nil
}

func (s *Service) GetDevice(ctx context.Context, id string) (device.Public, error) {
	nid, err := strictUUID(id)
	if err != nil {
		return device.Public{}, err
	}
	d, err := s.repo.Get(ctx, nid)
	if err != nil {
		return device.Public{}, err
	}
	return d.Public(), nil
}

func (s *Service) ListDevices(ctx context.Context) ([]device.Public, error) {
	devs, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]device.Public, 0, len(devs))
	for _, d := range devs {
		out = append(out, d.Public())
	}
	return out, nil
}

func (s *Service) SignTransaction(ctx context.Context, deviceID, data string) (signature, signedData string, err error) {
	nid, err := strictUUID(deviceID)
	if err != nil {
		return "", "", err
	}
	d, err := s.repo.Get(ctx, nid)
	if err != nil {
		return "", "", err
	}
	return d.Sign(data)
}
