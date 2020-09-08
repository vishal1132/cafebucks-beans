.Phony: protogo dev help

.DEFAULT_GOAL=help

help:
	@echo "protogo 				for compiling .proto file for go grpc plugin"
	@echo "dev					for running the server"
	@echo "buildImage   				for building docker image"

protogo:
	@protoc -I protos/ protos/beans.proto --go_out=plugins=grpc:protos/beans

dev:
	@go build -o beans ./server && PORT=8081 ./beans

buildImage:
	@docker build -t bean .