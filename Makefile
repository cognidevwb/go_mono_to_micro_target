# Workspace tasks. Every module builds on its own (GOWORK=off in CI) and together here.
MODULES := $(shell go list -m -f '{{.Dir}}')

.PHONY: sync build vet test up down hooks
sync:
	go work sync
build:
	@for m in $(MODULES); do (cd $$m && go build ./...) || exit 1; done
vet:
	@for m in $(MODULES); do (cd $$m && go vet ./...) || exit 1; done
test:
	@for m in $(MODULES); do (cd $$m && go test ./...) || exit 1; done
up:
	docker compose up --build
down:
	docker compose down -v
hooks:
	git config core.hooksPath .githooks
