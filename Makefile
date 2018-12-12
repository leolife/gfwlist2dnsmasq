
PREFIX := /usr/local

BINNAME := gfwlist2dnsmasq

default: build

test:
	gofmt -l *.go
	go vet
	go test -v cmd
build: 
	go build -o ${BINNAME}

clean:
	rm -f ${BINNAME}