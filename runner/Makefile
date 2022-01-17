BINARY_NAME=backend

build:
	go build -o ${BINARY_NAME} main.go


run:
	./${BINARY_NAME}

dev:
	go run main.go

build_and_run: build run

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod tidy

vet:
	go vet

lint:
	golangci-lint run --enable-all