# Include variables from the .envrc file
include .envrc

GO = go
VENDOR_BIN_PATH = ./bin/vendor

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
	@${GO} build -o ./bin/projectwall ./cmd/	