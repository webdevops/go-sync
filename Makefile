SOURCE = $(wildcard *.go)
TAG ?= $(shell git describe --tags)
GOBUILD = go build -ldflags '-w'

ALL = \
	$(foreach arch,x64 x32,\
	$(foreach suffix,linux osx windows,\
		build/gosync-$(suffix)-$(arch))) \
	$(foreach arch,arm arm64,\
		build/gosync-linux-$(arch))

all: test build

build: clean module test $(ALL)

# cram is a python app, so 'easy_install/pip install cram' to run tests
test:
	echo "No tests"
	#cram tests/*.test

clean:
	rm -f $(ALL)

module:
	go mod vendor

# os is determined as thus: if variable of suffix exists, it's taken, if not, then
# suffix itself is taken
osx = darwin
build/gosync-%-x64: $(SOURCE)
	@mkdir -p $(@D)
	CGO_ENABLED=0 GOOS=$(firstword $($*) $*) GOARCH=amd64 $(GOBUILD) -o $@

build/gosync-%-x32: $(SOURCE)
	@mkdir -p $(@D)
	CGO_ENABLED=0 GOOS=$(firstword $($*) $*) GOARCH=386 $(GOBUILD) -o $@

build/gosync-linux-arm: $(SOURCE)
	@mkdir -p $(@D)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 $(GOBUILD) -o $@

build/gosync-linux-arm64: $(SOURCE)
	@mkdir -p $(@D)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -o $@

release: build
	github-release release -u webdevops -r go-sync -t "$(TAG)" -n "$(TAG)" --description "$(TAG)"
	@for x in $(ALL); do \
		echo "Uploading $$x" && \
		github-release upload -u webdevops \
                              -r go-sync \
                              -t $(TAG) \
                              -f "$$x" \
                              -n "$$(basename $$x)"; \
	done
