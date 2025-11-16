##---------- Preliminaries ----------------------------------------------------
.POSIX:     # Get reliable POSIX behaviour
.SUFFIXES:  # Clear built-in inference rules

##---------- Variables --------------------------------------------------------
PREFIX = /usr/local  # Default installation directory
KO_DOCKER_REPO = sevaho/goforms

##---------- Build targets ----------------------------------------------------

##---------- Export .env as vars ----------------------------------------------
include .env
export

help: ## Show this help message (default)
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: help run css update lint test testr compose deploy serve release

run: ## Run application
	air

sqlgen: ## Generate SQL
	sqlc generate

sqlgenr: ## Generate SQL on repeat
	find db | entr -r make sqlgen

migrate: ## Run migrations
	go run . --migrate

css: ## Run CSS server
	tailwindcss -i ./src/web/assets/app.css -o ./src/web/static/css/app.css --watch

update: ## Update all dependencies
	go get -u
	go mod tidy

lint: ## Lint
	golangci-lint run --enable-all

test: ## Test
	ginkgo -r

testr: # Test with entr (rerun on file change)
	find . | entr -r ginkgo -r

compose: ## Run docker compose stack
	docker-compose rm -f
	docker-compose up

deploy: ## Deploy current build
	kubectx hetzner
	KO_DOCKER_REPO=$KO_DOCKER_REPO kubectl -n applications set image deployment/goforms goforms=$$(ko build .)

serve: ## Run docker locally
	KO_DOCKER_REPO=$KO_DOCKER_REPO docker run -p3000:3000 --network="host" --env-file=.env $$(ko build .) --serve

release: ## Release
	export KO_DOCKER_REPO=$(KO_DOCKER_REPO) && ko build --tags $$(git describe --tags --abbrev=0),latest --bare .
