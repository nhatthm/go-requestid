MODULE_NAME=requestid

VENDOR_DIR = vendor

GOLANGCI_LINT_VERSION ?= v2.1.6

GO ?= go
GOLANGCI_LINT ?= $(shell go env GOPATH)/bin/golangci-lint-$(GOLANGCI_LINT_VERSION)

.PHONY: $(VENDOR_DIR)
$(VENDOR_DIR):
	@mkdir -p $(VENDOR_DIR)
	@$(GO) mod vendor
	@$(GO) mod tidy

.PHONE: update
update:
	@$(GO) get -u ./...

.PHONY: tidy
tidy:
	@$(GO) mod tidy

.PHONY: lint
lint:
	@$(GOLANGCI_LINT) run

.PHONY: test
test: test-unit test-integration

## Run unit tests
.PHONY: test-unit
test-unit:
	@echo ">> unit test"
	@$(GO) test -gcflags=-l -coverprofile=unit.coverprofile -covermode=atomic -race ./...

.PHONY: test-integration
test-integration:
	@echo ">> integration test"
	@$(GO) test ./tests/integration/... -gcflags=-l -coverprofile=features.coverprofile -coverpkg ./... -race

.PHONY: $(GITHUB_OUTPUT)
$(GITHUB_OUTPUT):
	@echo "MODULE_NAME=$(MODULE_NAME)" >>"$@"
	@echo "GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION)" >>"$@"

$(GOLANGCI_LINT):
	@echo "$(OK_COLOR)==> Installing golangci-lint $(GOLANGCI_LINT_VERSION)$(NO_COLOR)"; \
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin "$(GOLANGCI_LINT_VERSION)"
	@mv ./bin/golangci-lint $(GOLANGCI_LINT)
