package mongo

import (
	"context"
	model "scrape-stock-market/models"
	"scrape-stock-market/stock"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StockRepository struct {
	db *mongo.Collection
}

func NewStockRepository(db *mongo.Database, collection string) stock.RepositoryImpl {
	return &StockRepository{
		db: db.Collection(collection),
	}
}

func (u *StockRepository) CreateOrUpdate(ctx context.Context, stock model.Stock) error {
	// find stock by symbol
	var stockData model.Stock
	u.db.FindOne(ctx, bson.M{"symbol": stock.Symbol}).Decode(&stockData)

	filter := bson.M{"symbol": stock.Symbol}

	changed := bson.M{}

	if stock.Symbol != "" {
		changed["symbol"] = stock.Symbol
	}

	if stock.Name != "" {
		changed["name"] = stock.Name
	}

	if stock.Index != "" {
		changed["index"] = stock.Index
	}

	if len(stock.PriceHistory) > 0 {
		stockData.PriceHistory = append(stockData.PriceHistory, stock.PriceHistory...)
		changed["price_history"] = stockData.PriceHistory
	}

	if len(stock.PriceHistoryDaily) > 0 {
		stockData.PriceHistoryDaily = append(stockData.PriceHistoryDaily, stock.PriceHistoryDaily...)
		changed["price_history_daily"] = stockData.PriceHistoryDaily
	}

	updateOptions := options.Update().SetUpsert(true)
	update := bson.M{"$set": changed}

	_, err := u.db.UpdateOne(ctx, filter, update, updateOptions)
	if err != nil {
		return err
	}

	return nil
}
