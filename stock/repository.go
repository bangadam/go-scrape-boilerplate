package stock

import (
	"context"
	model "scrape-stock-market/models"
)

type RepositoryImpl interface {
	CreateOrUpdate(ctx context.Context, stock model.Stock) error
}
