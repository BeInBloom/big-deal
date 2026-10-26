package maprepo

import (
	"context"
	"sync"

	"github.com/BeInBloom/big-deal/services/order/internal/models"
	"github.com/BeInBloom/big-deal/services/order/internal/services"
)

type MapRepo struct {
	mu     sync.RWMutex
	orders map[models.OrderId]models.OrderSnapshot
}

func New() *MapRepo {
	return &MapRepo{
		orders: make(map[models.OrderId]models.OrderSnapshot),
	}
}

func (r *MapRepo) Get(
	_ context.Context,
	id models.OrderId,
) (models.OrderSnapshot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[id]
	if !ok {
		return models.OrderSnapshot{}, services.ErrOrderNotFound
	}

	return cloneOrderSnapshot(order), nil
}

func (r *MapRepo) Create(
	_ context.Context,
	order models.OrderSnapshot,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.orders[order.Id]; ok {
		return services.ErrOrderAlreadyExists
	}

	r.orders[order.Id] = cloneOrderSnapshot(order)
	return nil
}

func (r *MapRepo) Update(
	_ context.Context,
	order models.OrderSnapshot,
	from models.OrderStatus,
) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	stored, ok := r.orders[order.Id]
	if !ok || stored.Status != from {
		return services.ErrOrderStatusConflict
	}

	r.orders[order.Id] = cloneOrderSnapshot(order)
	return nil
}

func cloneOrderSnapshot(order models.OrderSnapshot) models.OrderSnapshot {
	order.Parts = order.Parts.Clone()
	return order
}
