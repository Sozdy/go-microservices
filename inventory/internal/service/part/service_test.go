package part

import (
	"context"
	"testing"

	"github.com/Sozdy/go-microservices/inventory/internal/service/part/mocks"
)

type partFixture struct {
	ctx            context.Context
	partRepository *mocks.PartRepository
	service        *partService
}

func NewPartFixture(t *testing.T) *partFixture {
	t.Helper()

	partRepository := mocks.NewPartRepository(t)

	return &partFixture{
		ctx:            context.Background(),
		partRepository: partRepository,
		service:        NewPartService(partRepository),
	}
}
