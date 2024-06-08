# Include variables from the .envrc file
include .envrc

run:
	@go run ./src -port=${PORT} -database_url=${DATABASE_URL}