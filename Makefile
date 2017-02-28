OUTPUT_DIR=./_output

build: clean
	go build -o ${OUTPUT_DIR}/turbo-simulator ./cmd/turbo-simulator.go

test: clean build
	go test ./...

.PHONY: clean
clean:
	@: if [ -f ${OUTPUT_DIR} ] then rm -rf ${OUTPUT_DIR} fi
