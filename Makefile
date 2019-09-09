SOURCES := $(shell find . -name '*.go' -type f -not -path './vendor/*'  -not -path '*/mocks/*')

# Database
MDB_USER ?= 
MDB_PASSWORD ?= 
MDB_ADDRESS ?= 127.0.0.1:3306
MDB_DATABASE ?= employee

# Dependency Management
.PHONY: vendor
vendor: go.mod go.sum
	@GO111MODULE=on go get ./...

# Linter
.PHONY: lint-prepare
lint-prepare:
	@echo "Installing golangci-lint"
	@GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: lint
lint: vendor
	@echo "Start linting"
	@GO111MODULE=on golangci-lint-run ./..

# Mockery Prepare
.PHONY: mockery-prepare
mockery-prepare:
	@echo "Installing mockery"
	@GO111MODULE=off go get -u github.com/vektra/mockery/.../

# Database Migration
.PHONY: migrate-prepare
migrate-prepare:
	@echo "Prepare MariaDB migration"
	@GO111MODULE=off go get -tags 'mysql' -u github.com/golang-migrate/migrate/cmd/migrate

.PHONY: migrate-up
migrate-up:
	@echo "Start migrate up"
	@migrate -database "mysql://$(MDB_USER):$(MDB_PASSWORD)@tcp($(MDB_ADDRESS))/$(MDB_DATABASE)" \
	-path=driver/mariadb/migrations up

.PHONY: migrate-down
migrate-down:
	@echo "Start migrate down"
	@migrate -database "mysql://$(MDB_USER):$(MDB_PASSWORD)@tcp($(MDB_ADDRESS))/$(MDB_DATABASE)" \
	-path=driver/mariadb/migrations down

.PHONY: migrate-drop
migrate-drop:
	@echo "Start migrate drop"
	@migrate -database "mysql://$(MDB_USER):$(MDB_PASSWORD)@tcp($(MDB_ADDRESS))/$(MDB_DATABASE)" \
	-path=driver/mariadb/migrations drop

.PHONY: seed-up
seed-up:
	@echo "Start seed data"
	@migrate -database "mysql://$(MDB_USER):$(MDB_PASSWORD)@tcp($(MDB_ADDRESS))/$(MDB_DATABASE)" \
	-path=driver/mariadb/seeds up

.PHONY: seed-down
seed-down:
	@echo "Start unseed data"
	@migrate -database "mysql://$(MDB_USER):$(MDB_PASSWORD)@tcp($(MDB_ADDRESS))/$(MDB_DATABASE)" \
	-path=driver/mariadb/seeds down