package v1

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Sozdy/go-microservices/inventory/internal/errs"
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
		wantCode errs.Code
	}{
		{
			name:     "пустой uuid",
			request:  &inventoryv1.GetPartRequest{Uuid: ""},
			wantCode: errs.CodeInvalidArgument,
		},
		{
			name:     "невалидный uuid",
			request:  &inventoryv1.GetPartRequest{Uuid: "not-a-uuid"},
			wantCode: errs.CodeInvalidArgument,
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
			require.Equal(t, testCase.wantCode, errs.CodeOf(err))
		})
	}
}

func TestGetPart_ServiceError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		request    *inventoryv1.GetPartRequest
		serviceErr error
		wantCode   errs.Code
	}{
		{
			name:       "part не найдена",
			request:    &inventoryv1.GetPartRequest{Uuid: gofakeit.UUID()},
			serviceErr: errs.ErrPartNotFound,
			wantCode:   errs.CodeNotFound,
		},
		{
			name:       "невалидный uuid из сервиса",
			request:    &inventoryv1.GetPartRequest{Uuid: gofakeit.UUID()},
			serviceErr: errs.ErrInvalidUUID,
			wantCode:   errs.CodeInvalidArgument,
		},
		{
			name:       "внутренняя ошибка сервиса",
			request:    &inventoryv1.GetPartRequest{Uuid: gofakeit.UUID()},
			serviceErr: errors.New("что-то пошло не так в БД"),
			wantCode:   errs.CodeInternal,
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
			require.Equal(t, testCase.wantCode, errs.CodeOf(err))
		})
	}
}
