modules=./common ./common/errors ./common/containers ./common/volumes ./master ./master/client ./proxy ./nodeproxy

.PHONY: all
all: build

.PHONY: ensure
ensure:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

.PHONY: build
build: clean ensure bin/master bin/testtask bin/nodeproxy

.PHONY: test
test: ensure
	go test ${modules}

.PHONY: cover
cover: test
	go test -run '' -cover -coverprofile cover.out ${modules}
	go tool cover -html=cover.out

.PHONY: clean
clean:
	rm -f bin/master bin/testtask bin/nodeproxy

bin:
	mkdir -p bin

bin/master: bin
	go build -o bin/master master/*.go

bin/testtask: bin
	go build -o bin/testtask testtask/*.go

bin/nodeproxy: bin
	go build -o bin/nodeproxy nodeproxy/*.go

.PHONY: testimage
testimage:
	cd contrib/testimage && docker build -t creckx/testimage:latest .

.PHONY: lines
lines:
	@find ${modules} -name "*.go" -exec wc -l \{\} \; | cut -d " " -f 1 | xargs  |sed "s/ /+/g" | bc
