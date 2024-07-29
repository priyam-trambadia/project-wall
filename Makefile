# Include variables from the .envrc file
include .envrc

GO = go
VENDOR_BIN_PATH = ./bin/vendor
BUILD_OUTPUT = ./bin/projectwall

run: templ
	@${GO} run ./cmd -port=${PORT} -database_url=${DATABASE_URL} 

templ:
	@${VENDOR_BIN_PATH}/templ generate -path ./web/templates

.PHONEY: migrate
migrate:	
	@${VENDOR_BIN_PATH}/migrate -path=./migrations -database=${DATABASE_URL} down
	@${VENDOR_BIN_PATH}/migrate -path=./migrations -database=${DATABASE_URL} up	

.PHONEY: build
build:
	@${GO} build -o ${BUILD_OUTPUT} ./cmd/	

.PHONY: connect_db
connect_db:
	@psql ${DATABASE_URL}	

.PHONEY: install_vendor_bin
install_vendor_bin:
	@${GO} install -o ${VENDOR_BIN_PATH} github.com/a-h/templ/cmd/templ@v0.2.747
	@${GO} install -o ${VENDOR_BIN_PATH} github.com/golang-migrate/migrate/v4/cmd/migrate@latest	