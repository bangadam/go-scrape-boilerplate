package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"scrape-stock-market/stock"
	"sync"
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
	httpServer   *http.Server
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
	// defer scheduler.Stop()

	scheduler.AddFunc("@every 10s", a.scrapeData)
	go scheduler.Start()

	// trap SIGINT untuk trigger shutdown.
	select {}
}

func (a *App) scrapeData() {
	// Perform the scraping logic here
	// Use the stockUsecase to scrape the data and save it

	// Example:
	err := a.stockUsecase.ScrapeData()
	if err != nil {
		log.Printf("Failed to scrape data: %v", err)
	} else {
		log.Println("Data scraped successfully")
		// Do something with the scraped data, such as saving it to the database
	}
}

func (a *App) startHTTPServer(port string) error {
	// Initialize the HTTP server as before
	// ...

	// Start the HTTP server
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %v", err)
		}
	}()

	// Handle shutdown gracefully
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)
	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
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
		colly.Async(true),
		colly.UserAgent(viper.GetString("scrape.user_agent")),
	)

	return c
}
