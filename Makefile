cleandb:
	echo "" > ./db/sqlite3.db

migrate-up:
	if [ ! -d db ]; then mkdir db; fi
	chmod +x db/
	chmod +x db/sqlite3.db
	touch db/sqlite3.db
	sqlite3 db/sqlite3.db < migrate/account.sql
	sqlite3 db/sqlite3.db < migrate/block.sql
	sqlite3 db/sqlite3.db < migrate/tx.sql
	sqlite3 db/sqlite3.db < migrate/health.sql
	sqlite3 db/sqlite3.db < migrate/dummy.sql

migrate-down:
	echo "" > db/sqlite3.db

insert-dummy:
	sqlite3 db/sqlite3.db < migrate/dummy.sql

build:
	GOOS=linux GOARCH=amd64 go build -o ./bin/shimarin-chain *.go


