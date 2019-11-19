#!/usr/bin/env bash
go build -o main server.go
./main
#curl "localhost:8001/hello" -d "aruna2"