# Include variables from the .envrc file
include .envrc

BIN_PATH = ./bin

run: templ
	@go run ./src -port=${PORT} -database_url=${DATABASE_URL}

templ:
	@${BIN_PATH}/templ generate -path ./src/ui/