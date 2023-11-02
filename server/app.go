package server

import (
	"context"
	"encoding/json"
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

	stockRepo := _stockMongo.NewStockRepository(db, viper.GetString("mongo.collection"))

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

	scheduler.AddFunc("@every 10s", a.scrapeData)
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

	for _, stockSymbol := range stockSymbols {
		err := a.stockUsecase.ScrapeDataHistory(stockSymbol.Symbol)
		if err != nil {
			log.Printf("Failed to scrape data: %v", err)
			panic(err)
		} else {
			log.Println("Data scraped successfully on ", time.Now())
		}
	}
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
