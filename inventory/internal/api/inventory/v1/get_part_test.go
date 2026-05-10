package v1

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	inventoryErrs "github.com/Sozdy/go-microservices/inventory/internal/errors"
	"github.com/Sozdy/go-microservices/inventory/internal/model"
	inventoryv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/inventory/v1"
)

func TestGetPart_Success(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		request *inventoryv1.GetPartRequest
	}{
		{
			name: "Получение part",
			request: &inventoryv1.GetPartRequest{
				Uuid: gofakeit.UUID(),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// === Arrange ===
			fixture := newApiFixture(t)

			// === Expect ===
			fixture.service.EXPECT().
				GetPart(fixture.ctx, testCase.request.Uuid).
				Return(&model.Part{
					UUID: testCase.request.Uuid,
				}, nil).
				Once()

			// === Act ===
			partRes, err := fixture.api.GetPart(fixture.ctx, testCase.request)

			// === Assert ===
			assert.NoError(t, err)
			assert.Equal(t, testCase.request.Uuid, partRes.Part.Uuid)
		})
	}
}

func TestGetPart_ValidationError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		request  *inventoryv1.GetPartRequest
		wantCode codes.Code
	}{
		{
			name:     "пустой order_uuid",
			request:  &inventoryv1.GetPartRequest{Uuid: ""},
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "невалидный order_uuid",
			request:  &inventoryv1.GetPartRequest{Uuid: "not-a-uuid"},
			wantCode: codes.InvalidArgument,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// === Arrange ===
			fixture := newApiFixture(t)

			// === Expect ===

			// === Act ===
			response, err := fixture.api.GetPart(fixture.ctx, testCase.request)

			// === Assert ===
			require.Error(t, err)
			require.Nil(t, response)

			grpcStatus, isGRPCStatus := status.FromError(err)
			require.True(t, isGRPCStatus, "api должен возвращать grpc status error")
			require.Equal(t, testCase.wantCode, grpcStatus.Code())
		})
	}
}

func TestGetPart_ServiceError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		request    *inventoryv1.GetPartRequest
		serviceErr error
		wantCode   codes.Code
	}{
		{
			name:       "part не найдена",
			request:    &inventoryv1.GetPartRequest{Uuid: gofakeit.UUID()},
			serviceErr: inventoryErrs.ErrPartNotFound,
			wantCode:   codes.NotFound,
		},
		{
			name:       "невалидный uuid из сервиса",
			request:    &inventoryv1.GetPartRequest{Uuid: gofakeit.UUID()},
			serviceErr: inventoryErrs.ErrInvalidUUID,
			wantCode:   codes.InvalidArgument,
		},
		{
			name:       "внутренняя ошибка сервиса",
			request:    &inventoryv1.GetPartRequest{Uuid: gofakeit.UUID()},
			serviceErr: errors.New("что-то пошло не так в БД"),
			wantCode:   codes.Internal,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// === Arrange ===
			fixture := newApiFixture(t)

			// === Expect ===
			fixture.service.EXPECT().
				GetPart(fixture.ctx, testCase.request.Uuid).
				Return(nil, testCase.serviceErr).
				Once()

			// === Act ===
			response, err := fixture.api.GetPart(fixture.ctx, testCase.request)

			// === Assert ===
			require.Error(t, err)
			require.Nil(t, response)

			grpcStatus, isGRPCStatus := status.FromError(err)
			require.True(t, isGRPCStatus, "api должен возвращать grpc status error")
			require.Equal(t, testCase.wantCode, grpcStatus.Code())
		})
	}
}
