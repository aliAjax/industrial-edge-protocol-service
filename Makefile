GO=/tmp/go-toolchain/go/bin/go
fmt:
	$(GO)fmt -w $$(find . -name '*.go')
test:
	$(GO) test ./...
race:
	$(GO) test -race ./...
vet:
	$(GO) vet ./...
build:
	$(GO) build ./...
lines:
	find . -name '*.go' -not -name '*_test.go' -print0 | xargs -0 wc -l
smoke:
	./scripts/smoke.sh
