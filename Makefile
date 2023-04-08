#################### DEVELOPMENT ####################
run:
	go run src/main.go

dev-container-start:
	docker compose -f docker/docker-compose-dev.yml up

dev-container-stop:
	docker compose -f docker/docker-compose-dev.yml stop

mockery-repository:
	mockery --dir=src/domain --name=$(name) --filename=$(filename).go --output=src/domain/mocks/repository --outpkg=mocks

mockery-usecase:
	mockery --dir=src/domain --name=$(name) --filename=$(filename).go --output=src/domain/mocks/usecase --outpkg=mocks