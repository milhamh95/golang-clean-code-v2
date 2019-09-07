SOURCES := $(shell find . -name '*.go' -type f -not -path './vendor/*'  -not -path '*/mocks/*')

# Dependencies management
.PHONY: vendor
vendor: go.mod go.sum
	@GO111MODULE=on go get./...

# Linter
.PHONY: lint-prepare
lint-prepare:
	@echo "Installing golangci-lint"
	@GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: lint
lint: vendor
	@echo "Start linting"
	GO111MODULE=on golangci-lint-run ./...

# Database Migration
.PHONY: migrate-prepare
migrate-prepare:
	@echo: "Prepare MySQL migration"
	@GO111MODULE=off go get -tags 'mysql' -u github.com/golang-migrate/migrate/cmd/migrate

.PHONY: migrate-up
	@echo: "Start migrate up"
	@migrate -database "mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)" \
	-path=driver/mysql/migrations up

.PHONY: migrate-down
	@echo: "Start migrate down"
	@migrate -database "mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)" \
	-path=driver/mysql/migrations down

.PHONY: migrate-drop
	@echo: "Start migrate drop"
	@migrate -database "mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)" \
	-path=driver/mysql/migrations drop

.PHONY: seed-up
	@echo: "Start seed data"
	@migrate -database "mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)" \
	-path=driver/mysql/seeds up

.PHONY: seed-down
	@echo: "Start unseed data"
	@migrate -database "mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)" \
	-path=driver/mysql/seeds down