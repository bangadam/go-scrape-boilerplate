package usecase

import (
	"fmt"
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

func (u *StockUsecase) ScrapeDataHistory(symbol string) error {
	u.scrapeColly.OnHTML(`.BdT.Bdc\(\$seperatorColor\).Ta\(end\).Fz\(s\).Whs\(nw\)`, func(e *colly.HTMLElement) {
		date := e.ChildText(`td:nth-child(1)`)
		open := e.ChildText(`td:nth-child(2)`)
		high := e.ChildText(`td:nth-child(3)`)
		low := e.ChildText(`td:nth-child(4)`)
		close := e.ChildText(`td:nth-child(5)`)
		adjClose := e.ChildText(`td:nth-child(6)`)
		volume := e.ChildText(`td:nth-child(7)`)
		fmt.Println(date, open, high, low, close, adjClose, volume)
	})

	u.scrapeColly.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	u.scrapeColly.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	u.scrapeColly.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	u.scrapeColly.Visit(viper.GetString("scrape.url") + "/quote/" + symbol + "/history")

	return nil
}
