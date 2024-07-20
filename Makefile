# Include variables from the .envrc file
include .envrc

VENDOR_BIN_PATH = ./bin/vendor

run: templ
	@go run ./cmd -port=${PORT} -database_url=${DATABASE_URL} 

templ:
	@${VENDOR_BIN_PATH}/templ generate -path ./web/templates

.PHONEY: migrate
migrate:	
	@${VENDOR_BIN_PATH}/migrate -path=./migrations -database=${DATABASE_URL} down
	@${VENDOR_BIN_PATH}/migrate -path=./migrations -database=${DATABASE_URL} up	