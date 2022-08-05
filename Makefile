run:
	TRADER_LOG_LEVEL=4 \
	TRADER_DB_PATH=sqlite.db \
	TRADER_FREQUENCY=1000 \
	go run cmd/trader/main.go

docker:
	cd dev/ \
 	&& docker-compose up