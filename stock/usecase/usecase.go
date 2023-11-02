package usecase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	model "scrape-stock-market/models"
	"scrape-stock-market/stock"
	_repository "scrape-stock-market/stock/repository/mongo"

	"github.com/gocolly/colly"
	"github.com/spf13/viper"
)

type StockUsecase struct {
	stockRepo   *_repository.StockRepository
	scrapeColly *colly.Collector
}

func NewStockUsecase(stockrepo *_repository.StockRepository, scrapeColly *colly.Collector) stock.UsecaseImpl {
	return &StockUsecase{
		stockRepo:   stockrepo,
		scrapeColly: scrapeColly,
	}
}

func (u *StockUsecase) ScrapeData() error {
	// open file symbols.json from storage/symbols.json
	filePath := "./storage/symbols.json"
	symbols, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var stockSymbols []model.StockSymbol
	err = json.Unmarshal(symbols, &stockSymbols)
	if err != nil {
		return err
	}

	// scrape data from yahoo finance
	for _, symbol := range stockSymbols {
		// url := "https://finance.yahoo.com/quote/" + symbol.Symbol
		url := viper.GetString("scrape.url") + symbol.Link

		u.scrapeColly.OnHTML(`.Fw(b).Fz(36px).Mb(-4px).D(ib)`, func(e *colly.HTMLElement) {
			fmt.Println("Symbol:", symbol.Symbol)
		})

		// u.scrapeColly.OnRequest(func(r *colly.Request) {
		// 	fmt.Println("Visiting", r.URL)
		// })

		u.scrapeColly.Visit(url)
	}

	return nil
}
