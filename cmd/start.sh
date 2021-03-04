#!/bin/bash
cd ../api
protoc --go_out=plugins=grpc:. *.proto
cd ../cmd

cd account
go run . &
sleep 1s

cd ../article
go run . &
sleep 1s

cd ../comment
go run . &
sleep 1s

cd ../api-gateway
go run .

