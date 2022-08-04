run:
	TRADER_DB_PATH=sqlite.db \
	TRADER_FREQUENCY=100 \
	go run cmd/trader/main.go