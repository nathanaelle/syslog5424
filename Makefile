
export GOPATH:= $(CURDIR)/.GOPATH

test:
	go test -v ./...
