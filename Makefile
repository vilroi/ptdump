all: check build

check:
	go vet ./...

build:
	go build -o ptdump cmd/*.go
