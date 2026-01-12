#!/bin/bash
protoc --go_out=goclient \
  --go_opt=module=github.com/kinesis-network/go-greeter-client \
  --go-grpc_out=goclient \
  --go-grpc_opt=module=github.com/kinesis-network/go-greeter-client \
  grpcs/Protos/*.proto
cd goclient
go mod tidy
cd -
