COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

export GO111MODULE=on

.PHONY: build
build:
	@echo "-- building binary"
	go build \
		-ldflags "-X main.buildHash=${COMMIT} -X main.buildTime=${BUILD_TIME}" \
		-o ./bin/crudapp \
		./cmd/crudapp

.PHONY: test
test:
	@echo "-- run tests"
	go test -v -coverpkg=./... ./...

.PHONY: docker
docker:
	@echo "-- building docker container"
	docker build -f build/Dockerfile -t crudapp .

.PHONY: docker_multistage
docker_multistage: 
	@echo "-- building docker container - multistage"
	docker build -f build/Dockerfile.Multistage -t crudapp .

.PHONY: docker_run
docker_run: 
	@echo "-- starting docker container"
	docker run -it -p 8080:8080 crudapp

.PHONY: dc
dc:
	@echo "-- starting docker compose"
	docker-compose -f ./deployments/docker-compose.yml up

.PHONY: dcb
dcb:
	@echo "-- starting docker compose"
	docker-compose -f ./deployments/docker-compose.yml up --build
