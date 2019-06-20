INFO_COLOR=\033[1;34m
RESET=\033[0m
BOLD=\033[1m
dbcreate: ## Crate user and Create database
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Setup database$(RESET)"
	mysql -uroot -h127.0.0.1 -P3306 -e 'create database if not exists pdns DEFAULT CHARACTER SET ujis;'

dbdrop: ## Drop database
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Drop database$(RESET)"
	mysql -uroot -h127.0.0.1 -P3306 -e 'drop database pdns;'

dbmigrate: depends dbcreate
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Migrate database$(RESET)"
	sql-migrate up pdns
depends:
	GO111MODULE=off go get -v github.com/rubenv/sql-migrate/...

run:
	go run main.go --config ./misc/develop.toml server
