package v1

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	inventoryErrs "github.com/Sozdy/go-microservices/inventory/internal/errors"
	"github.com/Sozdy/go-microservices/inventory/internal/model"
	inventoryv1 "github.com/Sozdy/go-microservices/shared/pkg/proto/inventory/v1"
)

func TestListParts_Success(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name              string
		request           *inventoryv1.ListPartsRequest
		wantModelPartType model.PartType
	}{
		{
			name: "фильтр по типу UNSPECIFIED",
			request: &inventoryv1.ListPartsRequest{
				Uuids:    []string{gofakeit.UUID(), gofakeit.UUID()},
				PartType: inventoryv1.PartType_PART_TYPE_UNSPECIFIED,
			},
			wantModelPartType: model.PartTypeUnspecified,
		},
		{
			name: "фильтр по типу HULL",
			request: &inventoryv1.ListPartsRequest{
				Uuids:    []string{gofakeit.UUID(), gofakeit.UUID()},
				PartType: inventoryv1.PartType_PART_TYPE_HULL,
			},
			wantModelPartType: model.PartTypeHull,
		},
		{
			name: "фильтр по типу ENGINE",
			request: &inventoryv1.ListPartsRequest{
				Uuids:    []string{gofakeit.UUID(), gofakeit.UUID()},
				PartType: inventoryv1.PartType_PART_TYPE_ENGINE,
			},
			wantModelPartType: model.PartTypeEngine,
		},
		{
			name: "фильтр по типу SHIELD",
			request: &inventoryv1.ListPartsRequest{
				Uuids:    []string{gofakeit.UUID(), gofakeit.UUID()},
				PartType: inventoryv1.PartType_PART_TYPE_SHIELD,
			},
			wantModelPartType: model.PartTypeShield,
		},
		{
			name: "фильтр по типу WEAPON",
			request: &inventoryv1.ListPartsRequest{
				Uuids:    []string{gofakeit.UUID(), gofakeit.UUID()},
				PartType: inventoryv1.PartType_PART_TYPE_WEAPON,
			},
			wantModelPartType: model.PartTypeWeapon,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// === Arrange ===
			fixture := newApiFixture(t)

			// === Expect ===
			fixture.service.EXPECT().
				ListParts(fixture.ctx, testCase.request.Uuids, testCase.wantModelPartType).
				Return([]*model.Part{}, nil).
				Once()

			// === Act ===
			response, err := fixture.api.ListParts(fixture.ctx, testCase.request)

			// === Assert ===
			require.NoError(t, err)
			require.Equal(t, &inventoryv1.ListPartsResponse{}, response)
		})
	}
}

func TestListParts_ServiceError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name              string
		request           *inventoryv1.ListPartsRequest
		wantModelPartType model.PartType
		serviceErr        error
		wantCode          codes.Code
	}{
		{
			name: "одна из деталей не найдена",
			request: &inventoryv1.ListPartsRequest{
				Uuids:    []string{gofakeit.UUID(), gofakeit.UUID()},
				PartType: inventoryv1.PartType_PART_TYPE_HULL,
			},
			wantModelPartType: model.PartTypeHull,
			serviceErr:        inventoryErrs.ErrPartNotFound,
			wantCode:          codes.NotFound,
		},
		{
			name: "невалидный uuid из сервиса",
			request: &inventoryv1.ListPartsRequest{
				Uuids:    []string{gofakeit.UUID()},
				PartType: inventoryv1.PartType_PART_TYPE_ENGINE,
			},
			wantModelPartType: model.PartTypeEngine,
			serviceErr:        inventoryErrs.ErrInvalidUUID,
			wantCode:          codes.InvalidArgument,
		},
		{
			name: "внутренняя ошибка сервиса",
			request: &inventoryv1.ListPartsRequest{
				Uuids:    []string{gofakeit.UUID()},
				PartType: inventoryv1.PartType_PART_TYPE_WEAPON,
			},
			wantModelPartType: model.PartTypeWeapon,
			serviceErr:        errors.New("что-то пошло не так в БД"),
			wantCode:          codes.Internal,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// === Arrange ===
			fixture := newApiFixture(t)

			// === Expect ===
			fixture.service.EXPECT().
				ListParts(fixture.ctx, testCase.request.Uuids, testCase.wantModelPartType).
				Return(nil, testCase.serviceErr).
				Once()

			// === Act ===
			response, err := fixture.api.ListParts(fixture.ctx, testCase.request)

			// === Assert ===
			require.Error(t, err)
			require.Nil(t, response)

			grpcStatus, isGRPCStatus := status.FromError(err)
			require.True(t, isGRPCStatus, "api должен возвращать grpc status error")
			require.Equal(t, testCase.wantCode, grpcStatus.Code())
		})
	}
}
