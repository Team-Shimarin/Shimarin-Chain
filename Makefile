cleandb:
	rm -rf db/sqlite3.db

migrate-up:
	sqlite3 db/sqlite3.db < migrate/1518922917_accounts.up.sql
	sqlite3 db/sqlite3.db < migrate/1518948625_block.up.sql

migrate-down:
	sqlite3 db/sqlite3.db < migrate/1518922917_accounts.down.sql
	sqlite3 db/sqlite3.db < migrate/1518948625_block.down.sql

build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/anzu-chain *.go
