package part

import (
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	"github.com/Sozdy/go-microservices/inventory/internal/model"
)

func TestGetPart_Success(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		partUUID string
	}{
		{
			name:     "Получение part",
			partUUID: gofakeit.UUID(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// === Arrange ===
			fixture := NewPartFixture(t)

			// === Expect ===
			fixture.partRepository.EXPECT().
				GetPart(fixture.ctx, testCase.partUUID).
				Return(&model.Part{
					UUID: testCase.partUUID,
				}, nil).
				Once()

			// === Act ===
			part, err := fixture.service.GetPart(fixture.ctx, testCase.partUUID)

			// === Assert ===
			require.NoError(t, err)
			require.Equal(t, testCase.partUUID, part.UUID)
		})
	}
}

func TestGetPart_RepoError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		partUUID string
	}{
		{
			name:     "Получение ошибки при GetPart",
			partUUID: gofakeit.UUID(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// === Arrange ===
			fixture := NewPartFixture(t)

			// === Expect ===
			fixture.partRepository.EXPECT().
				GetPart(fixture.ctx, testCase.partUUID).
				Return(nil, errors.New("repo error")).
				Once()

			// === Act ===
			part, err := fixture.service.GetPart(fixture.ctx, testCase.partUUID)

			// === Assert ===
			require.Error(t, err)
			require.Nil(t, part)
		})
	}
}
