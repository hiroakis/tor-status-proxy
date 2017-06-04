APP_NAME = tor-status-proxy
IMAGE_NAME = tor-status-proxy
IMAGE_TAG = latest

lint:
	go vet ./...

build: lint
	GOOS=linux GOARCH=386 go build -ldflags "-s -w" -o $(APP_NAME)

run:
	go run *.go

clean:
	rm -f $(APP_NAME)

docker-build: build
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .

docker-run: docker-build
	docker run -d --name $(APP_NAME) -p 9000:9000 $(IMAGE_NAME):$(IMAGE_TAG)

docker-rm:
	docker rm -f $(shell docker ps -a | grep '$(APP_NAME)' | awk '{print $$1}')

docker-rmi:
	docker rmi -f $(shell docker images | grep '$(IMAGE_NAME)' | grep '$(IMAGE_TAG)' | awk '{print $$3}')

