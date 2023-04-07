#################### DEVELOPMENT ####################
run:
	go run src/main.go

dev-container-start:
	docker compose -f docker/docker-compose-dev.yml up

dev-container-stop:
	docker compose -f docker/docker-compose-dev.yml stop
