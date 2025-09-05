VERSION_PACKAGE := collector-go/internal/cmd/version

GO_LDFLAGS += -X $(VERSION_PACKAGE).hcgVersion=$(shell cat VERSION) \
	-X $(VERSION_PACKAGE).gitCommitID=$(GIT_COMMIT)

GIT_COMMIT:=$(shell git rev-parse HEAD)

.PHONY: dev
dev:
	go run cmd/main.go server --config etc/config.yml

.PHONY: build
build:
	CGO_ENABLED=0 go build -o bin -ldflags "$(GO_LDFLAGS)" cmd/main.go

.PHONY: prod
prod:
	bin/main server --config etc/config.yml
