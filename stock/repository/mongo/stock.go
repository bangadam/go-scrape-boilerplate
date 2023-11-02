package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type StockRepository struct {
	db *mongo.Collection
}

func NewStockRepository(db *mongo.Database, collection string) *StockRepository {
	return &StockRepository{
		db: db.Collection(collection),
	}
}
