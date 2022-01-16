#!/bin/bash

# download and install CompileDaemon without modifying go.mod
go install github.com/githubnemo/CompileDaemon@latest

# use ./main for command instead of "make run", make run exits after starting the application, and CompileDaemon will not have a reference to the application restart it on change
CompileDaemon --build="make build" --command="./elasticsearch-index-exporter -config example/config.yaml" --include=".env" --include="Makefile" --include="example/config.yaml"