package usecase

import (
	"context"
	"fmt"
	model "scrape-stock-market/models"
	"scrape-stock-market/stock"
	_repository "scrape-stock-market/stock"
	"scrape-stock-market/utils"
	"time"

	"github.com/gocolly/colly"
	"github.com/spf13/viper"
)

type StockUsecase struct {
	stockRepo   _repository.RepositoryImpl
	scrapeColly *colly.Collector
}

func NewStockUsecase(stockrepo _repository.RepositoryImpl, scrapeColly *colly.Collector) stock.UsecaseImpl {
	return &StockUsecase{
		stockRepo:   stockrepo,
		scrapeColly: scrapeColly,
	}
}

func (u *StockUsecase) ScrapeData(ctx context.Context, symbol string) error {
	var stock model.Stock

	// scrape profile
	profile := u.scrapeProfile(symbol)
	stock.Symbol = symbol
	stock.Profile = profile

	// scrape historical price data
	historicalPrice := u.scrapeHistoricalPrice(symbol)
	stock.HistoricalPrice = historicalPrice

	// scrape real price data
	realPrice := u.scrapeRealPrice(symbol)
	stock.RealPrice = realPrice

	// save to db
	err := u.stockRepo.CreateOrUpdate(ctx, stock)
	if err != nil {
		return err
	}

	return nil
}

func (u *StockUsecase) scrapeRealPrice(symbol string) (realPrice model.StockRealPrice) {
	newColly := u.scrapeColly.Clone()
	// Increase the timeout duration
	newColly.SetRequestTimeout(30 * time.Second) // Adjust the timeout value as needed
	var (
		currentPrice float64
		upDownPrice  string
	)

	newColly.OnHTML(`#quote-header-info > div.My\(6px\).Pos\(r\).smartphone_Mt\(6px\).W\(100\%\) > div.D\(ib\).Va\(m\).Maw\(65\%\).Ov\(h\) > div.D\(ib\).Mend\(20px\) > fin-streamer.Fw\(b\).Fz\(36px\).Mb\(-4px\).D\(ib\)`, func(e *colly.HTMLElement) {
		currentPrice = utils.StringToFloat64(e.Text)
	})

	newColly.OnHTML(`#quote-header-info > div.My\(6px\).Pos\(r\).smartphone_Mt\(6px\).W\(100\%\) > div.D\(ib\).Va\(m\).Maw\(65\%\).Ov\(h\) > div.D\(ib\).Mend\(20px\) > fin-streamer:nth-child(2) > span`, func(e *colly.HTMLElement) {
		upDownPrice = e.Text
	})

	newColly.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	newColly.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	newColly.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	// url https://https://finance.yahoo.com/quote/AAPL?p=AAPL
	newColly.Visit(viper.GetString("scrape.url") + "/quote/" + symbol + "?p=" + symbol)

	// set data
	realPrice.CurrentPrice = currentPrice
	realPrice.UpDownPrice = upDownPrice

	return realPrice
}

func (u *StockUsecase) scrapeHistoricalPrice(symbol string) (historicalPrices []model.StockHistoricalData) {
	newColly := u.scrapeColly.Clone()
	// Increase the timeout duration
	newColly.SetRequestTimeout(30 * time.Second) // Adjust the timeout value as needed

	newColly.OnHTML(`.BdT.Bdc\(\$seperatorColor\).Ta\(end\).Fz\(s\).Whs\(nw\)`, func(e *colly.HTMLElement) {
		date := e.ChildText(`td:nth-child(1)`)
		open := e.ChildText(`td:nth-child(2)`)
		high := e.ChildText(`td:nth-child(3)`)
		low := e.ChildText(`td:nth-child(4)`)
		close := e.ChildText(`td:nth-child(5)`)
		adjClose := e.ChildText(`td:nth-child(6)`)
		volume := e.ChildText(`td:nth-child(7)`)

		// parse
		dateParse := utils.StringDateToTime(date, "Jan 2, 2006")
		openParse := utils.StringToFloat64(open)
		highParse := utils.StringToFloat64(high)
		lowParse := utils.StringToFloat64(low)
		closeParse := utils.StringToFloat64(close)
		adjCloseParse := utils.StringToFloat64(adjClose)
		volumeParse := utils.StringToInt64(volume)

		historicalPrice := model.StockHistoricalData{
			Date:     dateParse,
			Open:     openParse,
			High:     highParse,
			Low:      lowParse,
			Close:    closeParse,
			AdjClose: adjCloseParse,
			Volume:   volumeParse,
		}

		historicalPrices = append(historicalPrices, historicalPrice)
	})

	newColly.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	newColly.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	newColly.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	newColly.Visit(viper.GetString("scrape.url") + "/quote/" + symbol + "/history?p=" + symbol)

	return historicalPrices
}

func (u *StockUsecase) scrapeProfile(symbol string) (profile model.StockProfile) {
	newColly := u.scrapeColly.Clone()
	// Increase the timeout duration
	newColly.SetRequestTimeout(30 * time.Second) // Adjust the timeout value as needed

	var company, address, sector, industry string

	// company
	newColly.OnHTML(`#Col1-0-Profile-Proxy > section > div.asset-profile-container > div > h3`, func(e *colly.HTMLElement) {
		company = e.Text
	})

	// address
	newColly.OnHTML(`#Col1-0-Profile-Proxy > section > div.asset-profile-container > div > div > p.D\(ib\).W\(47\.727\%\).Pend\(40px\)`, func(e *colly.HTMLElement) {
		address = e.Text
	})

	// sector
	newColly.OnHTML(`#Col1-0-Profile-Proxy > section > div.asset-profile-container > div > div > p.D\(ib\).Va\(t\) > span:nth-child(2)`, func(e *colly.HTMLElement) {
		sector = e.Text
	})

	// industry
	newColly.OnHTML(`#Col1-0-Profile-Proxy > section > div.asset-profile-container > div > div > p.D\(ib\).Va\(t\) > span:nth-child(5)`, func(e *colly.HTMLElement) {
		industry = e.Text
	})

	newColly.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	newColly.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// url https://https://finance.yahoo.com/quote/AAPL/profile?p=AAPL
	newColly.Visit(viper.GetString("scrape.url") + "/quote/" + symbol + "/profile?p=" + symbol)

	// set data
	profile.Company = company
	profile.Address = address
	profile.Sector = sector
	profile.Industry = industry

	return profile
}
