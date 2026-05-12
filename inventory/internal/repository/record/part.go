package record

type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         int64 // В копейках
	PartType      string
	StockQuantity int64
	CreatedAt     int64
}
