
.PHONY: clean all build run test-asm test-vmx 

usage:
	@echo "usage: make [all|build|clean|test]"

all: build

n2t-asm build: services/asm/*.go
	go build -o n2t-asm services/asm/main.go

clean:
	rm -f n2t coverage.out *.asm
	go clean -testcache

test: n2t
	go test -v ./pkg/asm
	go test -v ./pkg/vmx
	./n2t vmx pkg/vmx/testdata/SimpleAdd/SimpleAdd.vm

coverage:
	go test -race -covermode=atomic -coverprofile=coverage.out ./...

run: n2t
	./n2t asm pkg/asm/testdata/Max.asm

SimpleAdd.asm: n2t pkg/vmx/testdata/SimpleAdd/SimpleAdd.vm
	./n2t vmx pkg/vmx/testdata/SimpleAdd/SimpleAdd.vm
