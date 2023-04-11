#################### DEVELOPMENT ####################
## run: Running on Development
.PHONY: run
run:
	go run src/main.go

## swagger: Generate Swagger Docs
.PHONY: run
swagger:
	swag init -g src/main.go

dev-container-start:
	docker compose -f docker/docker-compose-dev.yml up

dev-container-stop:
	docker compose -f docker/docker-compose-dev.yml stop

mockery-repository:
	mockery --dir=src/domain --name=$(name) --filename=$(filename).go --output=src/domain/mocks/repository --outpkg=mocks

mockery-usecase:
	mockery --dir=src/domain --name=$(name) --filename=$(filename).go --output=src/domain/mocks/usecase --outpkg=mocks


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #
## test: Testing project 
.PHONY: test
test:
	@echo 'MyGramm on Test...'
	go test ./tests/...

## audit: Tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...

## vendor: Tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor