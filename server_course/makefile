BINARY_NAME=server_course

build:
	go build -o ${BINARY_NAME} cmd/server/main.go

run: build
	./${BINARY_NAME}

run_debug: build
	./${BINARY_NAME} --debug

clean:
	go clean
	rm ${BINARY_NAME}