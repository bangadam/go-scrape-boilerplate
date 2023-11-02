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