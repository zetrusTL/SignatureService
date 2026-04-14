package repository

import (
	"context"
	"sort"
	"sync"

	"github.com/zetrus/signature-service/internal/device"
	"github.com/zetrus/signature-service/internal/domain"
)

type Memory struct {
	mu   sync.RWMutex
	byID map[string]*device.Device
}

func NewMemory() *Memory {
	return &Memory{byID: make(map[string]*device.Device)}
}

func (m *Memory) Create(ctx context.Context, d *device.Device) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.byID[d.ID]; ok {
		return domain.ErrConflict
	}
	m.byID[d.ID] = d
	return nil
}

func (m *Memory) Get(ctx context.Context, id string) (*device.Device, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	d, ok := m.byID[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return d, nil
}

func (m *Memory) List(ctx context.Context) ([]*device.Device, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*device.Device, 0, len(m.byID))
	for _, d := range m.byID {
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].ID < out[j].ID
	})
	return out, nil
}
