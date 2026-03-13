.DEFAULT_GOAL := help

BINARY := openclaw-go-status
CMD := ./cmd/server

.PHONY: help
help:
	@printf "%s\n" \
	"Targets:" \
	"  make build        Build binary ($(BINARY))" \
	"  make run          Run server (go run $(CMD))" \
	"  make test         Run tests" \
	"  make test-cover   Run tests with coverage" \
	"  make fmt          Format code (gofmt -w)" \
	"  make fmt-check    Check formatting (gofmt -l)" \
	"  make vet          Run go vet" \
	"  make tidy         Run go mod tidy" \
	"  make clean        Remove build artifacts"

.PHONY: build
build:
	go build -o $(BINARY) $(CMD)

.PHONY: run
run:
	go run $(CMD)

.PHONY: test
test:
	go test ./...

.PHONY: test-cover
test-cover:
	go test -cover ./...

.PHONY: fmt
fmt:
	gofmt -w .

.PHONY: fmt-check
fmt-check:
	@test -z "$$(gofmt -l .)" || (gofmt -l . && exit 1)

.PHONY: vet
vet:
	go vet ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: clean
clean:
	rm -f $(BINARY)
