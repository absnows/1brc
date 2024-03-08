BINARY_NAME=brc

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin cmd/read/main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux cmd/read/main.go
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows cmd/read/main.go

run: 
	build
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows
