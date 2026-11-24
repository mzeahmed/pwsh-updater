BINARY  := pwsh-updater
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)
GO      ?= $(shell command -v go 2>/dev/null || echo /usr/local/go/bin/go)
BINDIR  ?= /usr/local/bin

.DEFAULT_GOAL := help
.PHONY: help build install uninstall run check fmt vet test clean

help: ## Affiche cette aide
	@grep -E '^[a-z-]+:.*##' $(MAKEFILE_LIST) | awk -F':.*## ' '{printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'

build: ## Compile le binaire dans ./pwsh-updater
	$(GO) build -ldflags "$(LDFLAGS)" -o $(BINARY) .

install: build ## Compile, installe dans /usr/local/bin (sudo) et retire le binaire local
	sudo install -m 755 $(BINARY) $(BINDIR)/$(BINARY)
	@rm -f $(BINARY)

uninstall: ## Supprime le binaire de /usr/local/bin (sudo)
	sudo rm -f $(BINDIR)/$(BINARY)

run: build ## Compile puis lance (ex: make run ARGS="--check")
	./$(BINARY) $(ARGS)

check: build ## Vérifie si une mise à jour de PowerShell existe
	./$(BINARY) --check

fmt: ## Formate le code
	$(GO) fmt ./...

vet: ## Analyse statique
	$(GO) vet ./...

test: ## Lance les tests
	$(GO) test ./...

clean: ## Supprime le binaire compilé
	rm -f $(BINARY)
