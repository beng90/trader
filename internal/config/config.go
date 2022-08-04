package config

type Config struct {
	Db        Db
	LogLevel  string
	Binance   Binance
	Frequency int
}

type Db struct {
	Path string
}

type Binance struct {
	ApiKey    string
	ApiSecret string
}
