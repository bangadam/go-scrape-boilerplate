package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	model "scrape-stock-market/models"
	"scrape-stock-market/stock"
	"sync"
	"syscall"
	"time"

	"github.com/gocolly/colly"
	cron "github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

	_stockMongo "scrape-stock-market/stock/repository/mongo"
	_stockUsecase "scrape-stock-market/stock/usecase"
)

type App struct {
	stockUsecase stock.UsecaseImpl
}

func NewApp() *App {
	db := initDB()
	initColly := initColly()

	// Initialize repository and usecase
	stockRepo := _stockMongo.NewStockRepository(db, model.StockCollection)

	return &App{
		stockUsecase: _stockUsecase.NewStockUsecase(stockRepo, initColly),
	}
}

func (a *App) Run(port string) error {
	// Start cron job
	a.startCronJob()

	// Start HTTP server
	// return a.startHTTPServer(port)
	return nil
}

func (a *App) startCronJob() {
	// set scheduler berdasarkan zona waktu sesuai kebutuhan
	jakartaTime, _ := time.LoadLocation("Asia/Jakarta")
	scheduler := cron.New(cron.WithLocation(jakartaTime))

	// stop scheduler tepat sebelum fungsi berakhir
	defer scheduler.Stop()

	// run cron job first time in a day
	scheduler.AddFunc("@every 10s", func() {
		fmt.Println("Running cron job at ", time.Now())
		a.scrapeData()
	})
	go scheduler.Start()

	// trap SIGINT untuk trigger shutdown.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

func (a *App) scrapeData() {
	// Perform the scraping logic here
	// open file symbols.json from storage/symbols.json
	filePath := "./storage/symbols.json"
	symbols, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		panic(err)
	}

	var stockSymbols []model.StockSymbol
	err = json.Unmarshal(symbols, &stockSymbols)
	if err != nil {
		log.Printf("Failed to unmarshal file: %v", err)
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(stockSymbols))

	// init context
	ctx := context.Background()

	for _, stockSymbol := range stockSymbols {
		go func(symbol string) {
			defer wg.Done()

			err := a.stockUsecase.ScrapeData(ctx, symbol)
			if err != nil {
				log.Printf("Failed to scrape data: %v", err)
			}
		}(stockSymbol.Symbol)
	}

	log.Println("Data scraped successfully on ", time.Now())

	// Wait for all goroutines to finish
	wg.Wait()
}

func initDB() *mongo.Database {
	var (
		mongoDatastore     *mongo.Database
		mongoDatastoreOnce sync.Once
	)

	mongoDatastoreOnce.Do(func() {
		databaseURI := viper.GetString("mongo.uri")
		connstr, err := connstring.Parse(databaseURI)
		if err != nil {
			log.Panicf("failed to parse mongo database uri %s", err.Error())
		}

		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connstr.String()))
		if err != nil {
			log.Panicf("failed to connect to mongo database %s", err.Error())
		}
		mongoDatastore = client.Database(connstr.Database)
	})
	return mongoDatastore
}

func initColly() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains(viper.GetString("scrape.domain")),
		colly.CacheDir(viper.GetString("scrape.cache")),
		colly.UserAgent(viper.GetString("scrape.user_agent")),
	)

	return c
}
