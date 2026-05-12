package v1

type api struct {
	orderService OrderService
}

func NewApi(orderService OrderService) *api {
	return &api{
		orderService: orderService,
	}
}
