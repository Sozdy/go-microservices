package errs

var (
	ErrPartNotFound = NotFound("деталь не найдена")
	ErrInvalidUUID  = InvalidArgument("неверный формат uuid")
)
