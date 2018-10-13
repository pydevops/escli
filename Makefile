SOURCEDIR=.
SOURCES :=  $(shell find . -name '*.go' -not -path './vendor/*')
BINARY=esops
VERSION=1.0.0
BUILD_TIME=`date +%FT%T%z`

PACKAGES := github.com/pydevops/esops
DEPENDENCIES := github.com/parnurzeal/gorequest


.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	go build  -o ${BINARY} ${SOURCES}

.PHONY: install
install:
	go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: test
test:
	go test -v $(PACKAGES)

.PHONY: silent-test
silent-test:
	go test $(PACKAGES)

.PHONY: format
format:
	go fmt $(PACKAGES)

.PHONY: deps
deps:
	go get $(DEPENDENCIES)
