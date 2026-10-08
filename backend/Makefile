.PHONY: gen lint tidy fmt go-mod-update test-unit test-integration test-component

gen:
	go generate ./internal/...
	$(MAKE) fmt

lint:
	go tool golangci-lint run ./...

tidy:
	go mod tidy

fmt:
	go tool gofumpt -l -w .

go-mod-update:
	go get -u ./...

test-unit:
	go test -race ./internal/...

test-integration:
	set -a; . ./.env.test; set +a; \
	go test -race -tags integration ./internal/... -count=1

test-component:
	set -a; . ./.env.test; set +a; \
	go test -race -p 1 -tags component ./tests/... -count=1
