INFO_COLOR=\033[1;34m
RESET=\033[0m
BOLD=\033[1m
TEST ?= $(shell go list ./... | grep -v -e vendor -e keys -e tmp)

BINARY:=pdns-api
SYSTEM:=
BUILDOPTS:=
GO ?= $(SYSTEM) go

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
	go install github.com/rubenv/sql-migrate/...@latest
	$(GO) get -u github.com/swaggo/echo-swagger

swag:
	swag init --parseInternal --parseDependency --parseDepth 2

run: swag
	$(GO) run main.go --config ./misc/develop.toml server
tidy:
	go mod tidy
test: tidy  ## Run test
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Testing$(RESET) (require: etcd,redis)"
	$(GO) test -v $(TEST) -timeout=30s -parallel=4
	$(GO) test -race $(TEST)

build_binary:
	$(GO) build $(BUILDOPTS) -ldflags="-s -w" -o $(BINARY)
