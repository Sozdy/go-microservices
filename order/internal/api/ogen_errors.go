package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ogen-go/ogen/ogenerrors"

	orderv1 "github.com/Sozdy/go-microservices/shared/pkg/openapi/order/v1"
)

// OgenErrorHandler перехватывает ошибки, которые ogen генерирует ДО вызова
// handler'ов: невалидный JSON, отсутствие required-полей, неверный формат UUID и т.д.
//
// Возвращает JSON в нашем формате orderv1.Error CreateOrderBadRequest/PayOrderNotFound/итд.,
// чтобы клиенты видели единый формат ошибок как при валидации от ogen, так и от domain.
func OgenErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	statusCode := http.StatusInternalServerError
	message := "внутренняя ошибка сервера"
	var details []orderv1.ErrorDetail

	var (
		decodeBodyErr    *ogenerrors.DecodeBodyError
		decodeParamsErr  *ogenerrors.DecodeParamsError
		decodeRequestErr *ogenerrors.DecodeRequestError
		securityErr      *ogenerrors.SecurityError
	)

	switch {
	case errors.As(err, &decodeBodyErr),
		errors.As(err, &decodeParamsErr),
		errors.As(err, &decodeRequestErr):
		statusCode = http.StatusBadRequest
		message = "невалидные параметры запроса"
		// Кпадём текст ошибки от ogen в details.
		// Field здесь "" точное поле ogen не всегда отдаёт структурировано.
		details = []orderv1.ErrorDetail{{
			Field:   "",
			Message: err.Error(),
		}}

	case errors.As(err, &securityErr):
		statusCode = http.StatusUnauthorized
		message = "требуется авторизация"

	default:
		slog.ErrorContext(ctx, "необработанная ошибка ogen", "err", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	body := orderv1.Error{
		Code:    statusCode,
		Message: message,
		Details: details,
	}
	if encErr := json.NewEncoder(w).Encode(body); encErr != nil {
		slog.ErrorContext(ctx, "не удалось сериализовать error response", "err", encErr)
	}
}
