modules=./common ./master ./proxy ./docker

.PHONY: all
all: build

.PHONY: ensure
ensure:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

.PHONY: build
build: clean ensure bin/master

.PHONY: test
test: ensure
	go test ${modules}

.PHONY: cover
cover: test
	go test -run '' -cover -coverprofile cover.out ${modules}
	go tool cover -html=cover.out

clean:
	rm bin/master

bin:
	mkdir -p bin

bin/master: bin
	go build -o bin/master master/main.go

