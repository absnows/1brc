BINARY_NAME=brc

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin cmd/main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux cmd/main.go
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows cmd/main.go

run: 
	build
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows

prepare:
	make gen
	make dev

dev:
	go run cmd/main.go -f=data/sample.txt

gen:
	go run cmd/produce/main.go -s=1000 -o=sample

test:
	go test ./...
