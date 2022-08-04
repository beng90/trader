package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/beng90/trader/internal/config"
	"github.com/beng90/trader/internal/models"
	"github.com/beng90/trader/internal/repositories"
	"github.com/beng90/trader/internal/services"
	"github.com/beng90/trader/pkg/logus"
	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var cfg config.Config

func init() {
	err := envconfig.Process("TRADER", &cfg)
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Start trader...")

	ctx := context.Background()
	l := log.New(os.Stdout, "", 5)
	stdLogger := logus.NewStdLogger(l)

	dbConfig := &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)}
	db, err := gorm.Open(sqlite.Open(cfg.Db.Path), dbConfig)
	checkErr(err)

	err = db.AutoMigrate(&models.Trade{}, &models.Order{})
	checkErr(err)

	client := binance.NewClient(cfg.Binance.ApiKey, cfg.Binance.ApiSecret)

	orderBookTickerRepository := repositories.NewOrderBookTickerRepository(client, stdLogger)
	tradeRepository := repositories.NewTradeRepository(db, stdLogger)
	orderRepository := repositories.NewOrderRepository(db, stdLogger)
	traderService := services.NewTraderService(stdLogger, orderBookTickerRepository, tradeRepository, orderRepository)

	for range time.Tick(time.Millisecond * time.Duration(cfg.Frequency)) {
		err = traderService.Watch(ctx)
		if err != nil {
			fmt.Errorf("%s\n", err)
		}
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
