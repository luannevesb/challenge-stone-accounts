export TERM=xterm-256color
export CLICOLOR_FORCE=true
export RICHGO_FORCE_COLOR=1

default: test

# Faz o Build
build: install
	@echo "Fazendo build..."
	env GOOS=linux go build -o bin/service cmd/main.go

test: install test-lint test-coverage

setup-local: install
	@go get -u golang.org/x/tools/...
	@go get -u golang.org/x/lint/golint
	@go get -u github.com/haya14busa/goverage
	@go get -u github.com/kyoh86/richgo
	@go get github.com/joho/godotenv/cmd/godotenv
	@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(GOPATH)/bin v1.15.0
	@golangci-lint --version

# Instala dependencias
install:
	@echo "Baixando depedencias..."
	@go mod verify

# Roda teste unitarios e gera coverage
test-coverage:
	@echo "Rodando testes"
	@richgo test -failfast -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html



# Testa o LINT com o GOLINT
test-lint:
	@golangci-lint run -c ./.golangci.yml ./...


clean:
	@go clean -modcache
	@rm -rf ./vendor




