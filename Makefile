
GOPATH:=$(shell go env GOPATH)
MODIFY=Mproto/imports/api.proto=github.com/micro/go-micro/v2/api/proto

.PHONY: proto
proto:
    
	protoc --proto_path=. --micro_out=${MODIFY}:. --go_out=${MODIFY}:. proto/WorkWeb/WorkWeb.proto
    

.PHONY: build
build: proto

	go build -o WorkWeb-service *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build -t ponywilliam/work-web .
	docker tag ponywilliam/work-web ponywilliam/work-web
	docker push ponywilliam/work-web