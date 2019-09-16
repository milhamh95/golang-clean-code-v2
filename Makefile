IMAGE_NAME := employee
TEST_OPTS := -covermode=atomic $(TEST_OPTS)

# Database
MYSQL_USER ?= employee
MYSQL_PASSWORD ?= employee-pass
MYSQL_ADDRESS ?= 127.0.0.1:3306
MYSQL_DATABASE ?= employee

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
	@migrate -database "mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)" \
	-path=driver/mariadb/migrations up

.PHONY: migrate-down
migrate-down:
	@echo "Start migrate down"
	@migrate -database "mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)" \
	-path=driver/mariadb/migrations down

.PHONY: migrate-drop
migrate-drop:
	@echo "Start migrate drop"
	@migrate -database "mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)" \
	-path=driver/mariadb/migrations drop

.PHONY: seed-up
seed-up:
	@echo "Start seed data"
	@migrate -database "mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)" \
	-path=driver/mariadb/seeds up

.PHONY: seed-down
seed-down:
	@echo "Start unseed data"
	@migrate -database "mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)" \
	-path=driver/mariadb/seeds down

# Docker
.PHONY: mariadb-up
mariadb-up:
	@echo "start mariadb"
	@docker-compose up -d mariadb

.PHONY: mariadb-down
mariadb-down:
	@echo "stop mariadb"
	@docker stop employee_mariadb

.PHONY: docker-dev
docker-dev:
	@echo "build employee DEV image"
	@docker build -t $(IMAGE_NAME) . -f Dockerfile.dev

.PHONY: docker
docker:
	@echo "build employee PROD image"
	@docker build -t $(IMAGE_NAME) .

.PHONY: run-dev
run-dev:
	@echo "run employee DEV"
	@docker-compose -f docker-compose.yaml -f docker-compose.dev.yaml up -d

.PHONY: run
run:
	@echo "run employee PROD"
	@docker-compose up -d

.PHONY: stop
stop:
	@docker-compose down

# Testing
.PHONY: unittest
unittest: vendor
	GO111MODULE=on go test -short $(TEST_OPTS) ./..

.PHONY: test
test: vendor
	GO111MODULE=on go test $(TEST_OPTS) ./...

# Mockery
DepartmentUseCase:
	@mockery -dir=domain -name=DepartmentUseCase -output=domain/mocks

DepartmentRepository:
	@mockery -dir=domain -name=DepartmentRepository -output=domain/mocks


