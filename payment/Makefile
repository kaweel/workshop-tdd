TZ ?= UTC
start:
	go run main.go
test:
	go test -v ./...
integration-test:
	go test -tags=integration_test ./...
unit-test:
	go test -tags=unit_test ./...
test-gen-cov:	
	go test -tags=unit_test,integration_test ./... -coverprofile=coverage.out
watch-cov-func:
	go tool cover -func=coverage.out
watch-cov:
	go tool cover -html=coverage.out
db-up:
	docker compose up
db-down:
	docker compose down