package model

import "time"

const StockCollection = "stock_yahoo"

type StockSymbol struct {
	Symbol string `json:"symbol" bson:"symbol"`
	Title  string `json:"title" bson:"title"`
	Link   string `json:"link" bson:"link"`
}

// expect struct is

type Stock struct {
	Symbol            string                   `bson:"symbol"`
	Name              string                   `bson:"name"`
	Index             string                   `bson:"index"`
	PriceHistory      []StockPriceHistory      `bson:"price_history"`
	PriceHistoryDaily []StockPriceHistoryDaily `bson:"price_history_daily"`
}

type StockPriceHistory struct {
	ScrapeTime       time.Time `bson:"scrape_time"`
	Price            string    `bson:"price"`
	UpDownPrice      string    `bson:"up_down_price"`
	UpDownPercentage string    `bson:"up_down_percentage"`
}

type StockPriceHistoryDaily struct {
	ScrapeTime time.Time `bson:"scrape_time"`
	Date       time.Time `json:"date" bson:"date"`
	Open       float64   `json:"open" bson:"open"`
	High       float64   `json:"high" bson:"high"`
	Low        float64   `json:"low" bson:"low"`
	Close      float64   `json:"close" bson:"close"`
	AdjClose   float64   `json:"adj_close" bson:"adj_close"`
	Volume     int64     `json:"volume" bson:"volume"`
}

type StockHistoricalData struct {
	Date     time.Time `json:"date" bson:"date"`
	Open     float64   `json:"open" bson:"open"`
	High     float64   `json:"high" bson:"high"`
	Low      float64   `json:"low" bson:"low"`
	Close    float64   `json:"close" bson:"close"`
	AdjClose float64   `json:"adj_close" bson:"adj_close"`
	Volume   int64     `json:"volume" bson:"volume"`
}

type StockProfile struct {
	Company  string `json:"company" bson:"company"`
	Address  string `json:"address" bson:"address"`
	Sector   string `json:"sector" bson:"sector"`
	Industry string `json:"industry" bson:"industry"`
}

type StockRealPrice struct {
	CurrentPrice float64 `json:"current_price" bson:"current_price"`
	UpDownPrice  string  `json:"up_down_price" bson:"up_down_price"`
}
