
# =============================================================================
run:
	go run cmd/apps/node/main.go

# =============================================================================
# Depends

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor