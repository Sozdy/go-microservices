package part

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	"github.com/Sozdy/go-microservices/inventory/internal/model"
)

func TestListParts_Success(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		partUUIDs []string
		partType  model.PartType
	}{
		{
			name:      "Получение parts",
			partUUIDs: []string{gofakeit.UUID(), gofakeit.UUID(), gofakeit.UUID()},
			partType:  model.PartTypeUnspecified,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// === Arrange ===
			fixture := NewPartFixture(t)

			// === Expect ===
			fixture.partRepository.EXPECT().
				ListParts(fixture.ctx, testCase.partUUIDs, testCase.partType).
				Return([]*model.Part{
					{UUID: testCase.partUUIDs[0]},
					{UUID: testCase.partUUIDs[1]},
					{UUID: testCase.partUUIDs[2]},
				}, nil).
				Once()

			// === Act ===
			parts, err := fixture.service.ListParts(fixture.ctx, testCase.partUUIDs, testCase.partType)

			// === Assert ===
			require.NoError(t, err)
			require.Equal(t, testCase.partUUIDs[0], parts[0].UUID)
			require.Equal(t, testCase.partUUIDs[1], parts[1].UUID)
			require.Equal(t, testCase.partUUIDs[2], parts[2].UUID)
		})
	}
}

func TestListParts_RepoError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		partUUIDs []string
		partType  model.PartType
	}{
		{
			name:      "Получение parts",
			partUUIDs: []string{gofakeit.UUID(), gofakeit.UUID(), gofakeit.UUID()},
			partType:  model.PartTypeUnspecified,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// === Arrange ===
			fixture := NewPartFixture(t)

			// === Expect ===
			ErrPartNotFound := errors.New("error")

			fixture.partRepository.EXPECT().
				ListParts(fixture.ctx, testCase.partUUIDs, testCase.partType).
				Return(nil, ErrPartNotFound).
				Once()

			// === Act ===
			parts, err := fixture.service.ListParts(fixture.ctx, testCase.partUUIDs, testCase.partType)

			// === Assert ===
			require.ErrorIs(t, err, ErrPartNotFound)
			require.Nil(t, parts)
		})
	}
}
