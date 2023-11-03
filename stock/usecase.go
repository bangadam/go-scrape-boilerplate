package stock

import "context"

type UsecaseImpl interface {
	ScrapeData(ctx context.Context, symbol string) error
}
