BINARY_NAME=elasticsearch-index-exporter
MAIN_FILE=elasticsearch-index-exporter.go

.PHONY: build
build:
	go build -o ${BINARY_NAME} ${MAIN_FILE}

.PHONY: run
run:
	./${BINARY_NAME}

.PHONY: clean
clean:
	go clean
	rm -f ${BINARY_NAME}