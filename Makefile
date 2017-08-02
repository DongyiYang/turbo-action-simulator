OUTPUT_DIR=./_output

product: clean
	env GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_DIR}/turbo-simulator.linux ./cmd/turbo-simulator.go

build: clean
	go build -o ${OUTPUT_DIR}/turbo-simulator ./cmd/turbo-simulator.go

test: clean build
	go test ./...

.PHONY: clean
clean:
	@: if [ -f ${OUTPUT_DIR} ] then rm -rf ${OUTPUT_DIR} fi
