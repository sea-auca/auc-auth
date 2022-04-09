#docker image name
IMAGE_NAME = sea.auca.kg/auc-auth
IMAGE_VERSION = 0.0.1

.SILENT: run clean

include .env
export

# DEVELOPMENT OPERATIONS

build: mod_tidy create_docker tag_latest
	
run:
	go run cmd/cmd.go

#Run docker image on host network
drun: down_container image
	docker run --name $(IMAGE_NAME) --network=host -d $(IMAGE_NAME):latest
# PRODUCTION BUILDS


create_docker:
	docker build --tag $(IMAGE_NAME):$(IMAGE_VERSION) .

tag_latest:
	docker image tag $(IMAGE_NAME):$(IMAGE_VERSION) $(IMAGE_NAME):latest

mod_tidy: #preparation step - in case of go.sum file is missing
	go mod tidy

down_container:
	docker stop $(IMAGE_NAME) || true && docker rm $(IMAGE_NAME) || true

# CLEANING

clean: 
	rm -f bin/main