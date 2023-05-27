METADATA_BINARY=metadataApp
DEAL_BINARY=dealApp
RATING_BINARY=ratingApp

start: build_metadata build_deal build_rating
	docker-compose down
	@echo "Docker images closed!"
	docker-compose up --build -d 
	@echo "Docker images built and started!"

stop:
	docker-compose down
	@echo "Docker images closed!"

build: build_metadata build_deal build_rating
	@echo "Docker images built!"

build_metadata:
	@echo "Building metadata binary..."
	cd ./metadata && env GOOS=linux CGO_ENABLED=0 go build -o ${METADATA_BINARY} ./
	@echo "Done!"

build_deal:
	@echo "Building deal binary..."
	cd ./deal && env GOOS=linux CGO_ENABLED=0 go build -o ${DEAL_BINARY} ./
	@echo "Done!"

build_rating:
	@echo "Building rating binary..."
	cd ./rating && env GOOS=linux CGO_ENABLED=0 go build -o ${RATING_BINARY} ./
	@echo "Done!"
