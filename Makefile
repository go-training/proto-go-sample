GO ?= go

.PHONY: build
build: server chi gin client

.PHONY: server
server: gin chi

.PHONY: chi
chi:
	$(GO) build -o bin/$@-server cmd/server/$@/*.go

.PHONY: gin
gin:
	$(GO) build -o bin/$@-server cmd/server/$@/*.go

.PHONY: client
client:
	$(GO) build -o bin/$@ cmd/$@/main.go

.PHONY: upgrade
upgrade: ## Upgrade dependencies
	$(GO) get -u -t ./... && go mod tidy -v

test:
	@$(GO) test -cover ./... && echo "\n==>\033[32m Ok\033[m\n" || exit 1

clean:
	rm -rf gen bin
