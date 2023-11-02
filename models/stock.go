package model

type StockSymbol struct {
	Symbol string `json:"symbol" bson:"symbol"`
	Title  string `json:"title" bson:"title"`
	Link   string `json:"link" bson:"link"`
}

type Stock struct {
	ID        string `bson:"_id"`
	Symbol    string `bson:"symbol"`
	Title     string `bson:"title"`
	URLDetail string `bson:"url_detail"`
}

type HistoricalData struct {
	Date     string `json:"date" bson:"date"`
	Open     string `json:"open" bson:"open"`
	High     string `json:"high" bson:"high"`
	Low      string `json:"low" bson:"low"`
	Close    string `json:"close" bson:"close"`
	AdjClose string `json:"adj_close" bson:"adj_close"`
	Volume   string `json:"volume" bson:"volume"`
}
