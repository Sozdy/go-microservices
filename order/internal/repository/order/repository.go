package order

import (
	"sync"

	"github.com/google/uuid"

	"github.com/Sozdy/go-microservices/order/internal/repository/record"
)

type repo struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]record.Order
}

func NewRepository() *repo {
	return &repo{
		orders: make(map[uuid.UUID]record.Order),
	}
}
