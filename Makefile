PACKAGE  = fhirbase
export GOPATH   = $(CURDIR)/.gopath
BASE     = $(GOPATH)/src/$(PACKAGE)
DATE    ?= $(shell date +%FT%T%z)
VERSION ?= $(shell (cat $(BASE)/.version 2> /dev/null) || (echo 'nightly-\c' && git rev-parse --short HEAD 2> /dev/null)  | tr -d "\n")

GO      = go
GODOC   = godoc
GOFMT   = gofmt


.PHONY: all
all: lint fmt | $(BASE)
	$(GO) build \
	-v \
	-tags release \
	-ldflags '-X "main.Version=$(VERSION)" -X "main.BuildDate=$(DATE)"' \
	-o bin/$(PACKAGE)$(BINSUFFIX) *.go

$(BASE):
	@mkdir -p $(dir $@)
	@ln -sf $(CURDIR) $@

# Tools

.PHONY: lint
lint: $(BASE) $(GOLINT)
	$Q cd $(BASE) && ret=0 && for pkg in $(PKGS); do \
	test -z "$$($(GOLINT) $$pkg | tee /dev/stderr)" || ret=1 ; \
	done ; exit $$ret

.PHONY: fmt
fmt:
	@ret=0 && for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
	$(GOFMT) -l -w $$d/*.go || ret=$$? ; \
	done ; exit $$ret

.PHONY: clean
clean:
	go clean -modcache
	rm -rf bin .gopath vendor

.PHONY: tests
test: fmt lint
	go test $(ARGS)

# .PHONY: docker
# docker: Dockerfile bin/fhirbase-linux-amd64
# 	docker build . -t fhirbase/fhirbase:$(VERSION) -t fhirbase/fhirbase:latest && \
# 	docker push fhirbase/fhirbase:$(VERSION) && \
# 	docker push fhirbase/fhirbase:latest
