
all: clean test build

test:
	go test ./...

build:
	go build -o bin/vwap main.go

run:
	go run main.go

serve:
	RUNNING_MODE=server go run main.go

clean:
	go clean
	rm -rf bin/
