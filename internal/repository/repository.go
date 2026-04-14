package repository

import (
	"context"

	"github.com/zetrus/signature-service/internal/device"
)

type DeviceRepository interface {
	Create(ctx context.Context, d *device.Device) error
	Get(ctx context.Context, id string) (*device.Device, error)
	List(ctx context.Context) ([]*device.Device, error)
}
