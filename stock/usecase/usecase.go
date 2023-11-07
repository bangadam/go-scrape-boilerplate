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

func NewStockUsecase(stockrepo _repository.RepositoryImpl) stock.UsecaseImpl {
	return &StockUsecase{
		stockRepo: stockrepo,
	}
}

func initColly() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains(viper.GetString("scrape.domain")),
		colly.CacheDir(viper.GetString("scrape.cache")),
		colly.UserAgent(viper.GetString("scrape.user_agent")),
	)

	return c
}

func (u *StockUsecase) ScrapeData(ctx context.Context, symbol string) error {
	var stock model.Stock
	scrapeTime := time.Now()

	// init colly
	u.scrapeColly = initColly()

	// scrape profile
	profile := u.scrapeProfile(symbol)
	stock.Symbol = symbol
	stock.Name = profile.Company

	// scrape historical price data
	stockPriceHistoryDaily := u.scrapeHistoricalPrice(symbol, scrapeTime)
	stock.PriceHistoryDaily = stockPriceHistoryDaily

	// scrape real price data
	realPrice := u.scrapeRealPrice(symbol, scrapeTime)
	stock.PriceHistory = append(stock.PriceHistory, realPrice)

	// convert to json
	// realPriceJSON, err := json.Marshal(stock)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(realPriceJSON))
	// save to db
	err := u.stockRepo.CreateOrUpdate(ctx, stock)
	if err != nil {
		return err
	}

	return nil
}

func (u *StockUsecase) scrapeRealPrice(symbol string, scrapeTime time.Time) (realPrice model.StockPriceHistory) {
	newColly := u.scrapeColly.Clone()
	// Increase the timeout duration
	// newColly.SetRequestTimeout(30 * time.Second) // Adjust the timeout value as needed
	var (
		currentPrice     string
		upDownPrice      string
		upDownPercentage string
	)

	newColly.OnHTML(`#quote-header-info > div.My\(6px\).Pos\(r\).smartphone_Mt\(6px\).W\(100\%\) > div.D\(ib\).Va\(m\).Maw\(65\%\).Ov\(h\) > div > fin-streamer.Fw\(b\).Fz\(36px\).Mb\(-4px\).D\(ib\)`, func(e *colly.HTMLElement) {
		currentPrice = e.Attr("value")
	})

	newColly.OnHTML(`#quote-header-info > div.My\(6px\).Pos\(r\).smartphone_Mt\(6px\).W\(100\%\) > div.D\(ib\).Va\(m\).Maw\(65\%\).Ov\(h\) > div.D\(ib\).Mend\(20px\) > fin-streamer:nth-child(2) > span`, func(e *colly.HTMLElement) {
		upDownPrice = e.Text
	})

	newColly.OnHTML(`#quote-header-info > div.My\(6px\).Pos\(r\).smartphone_Mt\(6px\).W\(100\%\) > div.D\(ib\).Va\(m\).Maw\(65\%\).Ov\(h\) > div.D\(ib\).Mend\(20px\) > fin-streamer:nth-child(3) > span`, func(e *colly.HTMLElement) {
		upDownPercentage = e.Text
	})

	newColly.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// url https://https://finance.yahoo.com/quote/AAPL?p=AAPL
	newColly.Visit(viper.GetString("scrape.url") + "/quote/" + symbol + "?p=" + symbol)

	// set data
	realPrice.ScrapeTime = scrapeTime
	realPrice.Price = currentPrice
	realPrice.UpDownPrice = upDownPrice
	realPrice.UpDownPercentage = upDownPercentage

	return realPrice
}

func (u *StockUsecase) scrapeHistoricalPrice(symbol string, scrapeTime time.Time) (historicalPrices []model.StockPriceHistoryDaily) {
	newColly := u.scrapeColly.Clone()
	// Increase the timeout duration
	// newColly.SetRequestTimeout(30 * time.Second) // Adjust the timeout value as needed

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

		stockPriceHistoryDaily := model.StockPriceHistoryDaily{
			ScrapeTime: scrapeTime,
			Date:       dateParse,
			Open:       openParse,
			High:       highParse,
			Low:        lowParse,
			Close:      closeParse,
			AdjClose:   adjCloseParse,
			Volume:     volumeParse,
		}

		historicalPrices = append(historicalPrices, stockPriceHistoryDaily)
	})

	newColly.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	newColly.Visit(viper.GetString("scrape.url") + "/quote/" + symbol + "/history?p=" + symbol)

	return historicalPrices
}

func (u *StockUsecase) scrapeProfile(symbol string) (profile model.StockProfile) {
	newColly := u.scrapeColly.Clone()
	// Increase the timeout duration
	// newColly.SetRequestTimeout(30 * time.Second) // Adjust the timeout value as needed

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
