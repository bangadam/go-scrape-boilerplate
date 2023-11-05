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

	if !stock.Data.YahooFinance.ScrapeTime.IsZero() {
		changed["data.yahoo_finance.scrape_time"] = stock.Data.YahooFinance.ScrapeTime
	}

	if stock.Data.YahooFinance.Profile.Company != "" {
		changed["data.yahoo_finance.profile.company"] = stock.Data.YahooFinance.Profile.Company
	}
	if stock.Data.YahooFinance.Profile.Address != "" {
		changed["data.yahoo_finance.profile.address"] = stock.Data.YahooFinance.Profile.Address
	}
	if stock.Data.YahooFinance.Profile.Sector != "" {
		changed["data.yahoo_finance.profile.sector"] = stock.Data.YahooFinance.Profile.Sector
	}
	if stock.Data.YahooFinance.Profile.Industry != "" {
		changed["data.yahoo_finance.profile.industry"] = stock.Data.YahooFinance.Profile.Industry
	}

	if len(stock.Data.YahooFinance.HistoricalPrice) > 0 {
		changed["data.yahoo_finance.historical_price"] = stock.Data.YahooFinance.HistoricalPrice
	}

	if stock.Data.YahooFinance.RealPrice.CurrentPrice != 0 {
		changed["data.yahoo_finance.real_price.current_price"] = stock.Data.YahooFinance.RealPrice.CurrentPrice
	}
	if stock.Data.YahooFinance.RealPrice.UpDownPrice != "" {
		changed["data.yahoo_finance.real_price.up_down_price"] = stock.Data.YahooFinance.RealPrice.UpDownPrice
	}

	updateOptions := options.Update().SetUpsert(true)
	update := bson.M{"$set": changed}

	_, err := u.db.UpdateOne(ctx, filter, update, updateOptions)
	if err != nil {
		return err
	}

	return nil
}
