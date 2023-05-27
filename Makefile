METADATA_BINARY=metadataApp

start: build_metadata
	docker-compose down
	@echo "Docker images closed!"
	docker-compose up --build -d 
	@echo "Docker images built and started!"

stop:
	docker-compose down
	@echo "Docker images closed!"

build: build_metadata
	@echo "Docker images built!"

build_metadata:
	@echo "Building metadata binary..."
	cd ./metadata && env GOOS=linux CGO_ENABLED=0 go build -o ${METADATA_BINARY} ./
	@echo "Done!"
