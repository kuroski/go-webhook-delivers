include .env

# ==================================================================================== #
# CONSTANTS
# ==================================================================================== #
IMAGE_NAME=go-webhook-deliveries
CONTAINER_NAME=go-webhook-deliveries

# ==================================================================================== #
# HELPERS
# ==================================================================================== #
## Prints usage information for each target
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #
## Docker Run: Runs Docker container for the web server
.PHONY: web-run-docker
web-run-docker:
	docker build -t $(IMAGE_NAME) .
	docker rm -f $(CONTAINER_NAME) || true
	docker run -d -p 8080:80 --env-file .env --name $(CONTAINER_NAME) $(IMAGE_NAME)
	@echo 'Your server is running on localhost:8080'

## Web Application: Runs the web application without docker
.PHONY: web-run
web-run:
	go run ./cmd/web

## Client: Runs the client CLI
.PHONY: cli-run
cli-run:
	@source_arg=$${source:-${DEV_CLI_SOURCE_URL}}; \
	target_arg=$${target:-${DEV_CLI_TARGET_URL}}; \
	go run ./cmd/cli --source=$$source_arg --target=$$target_arg

## Dev: Runs dev environment with docker compose
.PHONY: dev
dev:
	docker compose -f compose.dev.yml up
	@echo 'Your server is running on localhost:3000'


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #
## Tidy: Formats .go files and tidies dependencies
.PHONY: tidy
tidy:
	@echo 'Formatting .go files...'
	go fmt ./...
	@echo 'Tidying module dependencies...'
	go mod tidy

## Audit: Runs quality control checks
.PHONY: audit
audit:
	@echo 'Checking module dependencies'
	go mod tidy -diff
	go mod verify
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...


# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #
## Deploy the application (skipping Kamal hooks)
.PHONY: deploy
deploy:
	@echo 'Deploying the app...'
	kamal deploy -H
	@echo 'Tidying module dependencies...'
	go mod tidy