all: check build

check:
	go vet ./...

build:
	go build -o ptdump cmd/*.go

clean:
	rm ptdump

.PHONY: clean
