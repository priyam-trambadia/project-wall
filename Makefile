# Include variables from the .envrc file
include .envrc

BIN_PATH = ./bin

run: templ
	@go run ./cmd -port=${PORT} -database_url=${DATABASE_URL} 

templ:
	@${BIN_PATH}/templ generate -path ./web/templates