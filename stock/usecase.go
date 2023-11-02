package stock

type UsecaseImpl interface {
	ScrapeDataHistory(symbol string) error
}
