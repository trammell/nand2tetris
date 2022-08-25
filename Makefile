
.PHONY: clean all build run test-asm test-vmx

usage:
	@echo "usage: make [all|build|clean|test|tools]"

all: build

n2t build: cmd/*.go pkg/asm/*.go pkg/vmx/*.go
	go build -o n2t main.go

clean:
	rm -f n2t coverage.out *.asm
	go clean -testcache ./...

realclean:
	rm -rf tools/

test-asm:
	go test -v ./pkg/asm

coverage:
	go test -race -covermode=atomic -coverprofile=coverage.out ./...

run: n2t
	./n2t asm pkg/asm/testdata/Max.asm

test-vmx: n2t
	go test -v ./pkg/vmx
	./n2t vmx pkg/vmx/testdata/SimpleAdd/SimpleAdd.vm


SimpleAdd.asm: n2t pkg/vmx/testdata/SimpleAdd/SimpleAdd.vm
	./n2t vmx pkg/vmx/testdata/SimpleAdd/SimpleAdd.vm

tools: $(HOME)/Downloads/nand2tetris.zip
	unzip -o $(HOME)/Downloads/nand2tetris.zip "nand2tetris/tools/*" -d tools

# special commands to create the Coursera submission file
project7.zip:
	zip project7.zip blargh

