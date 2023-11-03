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
	filter := bson.M{"symbol": stock.Symbol}

	changed := bson.M{}

	if stock.Profile.Company != "" {
		changed["profile.company"] = stock.Profile.Company
	}
	if stock.Profile.Address != "" {
		changed["profile.address"] = stock.Profile.Address
	}
	if stock.Profile.Sector != "" {
		changed["profile.sector"] = stock.Profile.Sector
	}
	if stock.Profile.Industry != "" {
		changed["profile.industry"] = stock.Profile.Industry
	}

	if len(stock.HistoricalPrice) > 0 {
		changed["historical_price"] = stock.HistoricalPrice
	}

	if stock.RealPrice.CurrentPrice != 0 {
		changed["real_price.current_price"] = stock.RealPrice.CurrentPrice
	}
	if stock.RealPrice.UpDownPrice != "" {
		changed["real_price.up_down_price"] = stock.RealPrice.UpDownPrice
	}

	updateOptions := options.Update().SetUpsert(true)
	update := bson.M{"$set": changed}

	_, err := u.db.UpdateOne(ctx, filter, update, updateOptions)
	if err != nil {
		return err
	}

	return nil
}
