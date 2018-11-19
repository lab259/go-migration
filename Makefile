VERSION ?= $(shell git describe --tags 2>/dev/null | cut -c 2-)
TEST_FLAGS ?=
REPO_OWNER ?= $(shell cd .. && basename "$$(pwd)")

GOPATH=$(CURDIR)/../../../../
GOPATHCMD=GOPATH=$(GOPATH)

COVERDIR=$(CURDIR)/.cover
COVERAGEFILE=$(COVERDIR)/cover.out

test:
	@${GOPATHCMD} ginkgo --failFast ./...

test-watch:
	@${GOPATHCMD} ginkgo watch -cover -r ./...

coverage:
	@mkdir -p $(COVERDIR)
	@${GOPATHCMD} ginkgo -r -covermode=count --cover --trace ./
	@echo "mode: count" > "${COVERAGEFILE}"
	@find . -type f -name *.coverprofile -exec grep -h -v "^mode:" {} >> "${COVERAGEFILE}" \; -exec rm -f {} \;

coverage-ci:
	@mkdir -p $(COVERDIR)
	@${GOPATHCMD} ginkgo -r -covermode=count --cover --trace ./
	@echo "mode: count" > "${COVERAGEFILE}"
	@find . -type f -name *.coverprofile -exec grep -h -v "^mode:" {} >> "${COVERAGEFILE}" \; -exec rm -f {} \;

coverage-html:
	@$(GOPATHCMD) go tool cover -html="${COVERAGEFILE}" -o .cover/report.html

dep-ensure:
	@$(GOPATHCMD) dep ensure -v

dep-update:
	@$(GOPATHCMD) dep ensure -v

vet:
	@$(GOPATHCMD) go vet ./...

lint:
	@$(GOPATHCMD) golint

fmt:
	@$(GOPATHCMD) go fmt ./...

# example: fswatch -0 --exclude .godoc.pid --event Updated . | xargs -0 -n1 -I{} make docs
docs:
	-make kill-docs
	nohup godoc -play -http=127.0.0.1:6064 </dev/null >/dev/null 2>&1 & echo $$! > .godoc.pid
	cat .godoc.pid

kill-docs:
	@cat .godoc.pid
	kill -9 $$(cat .godoc.pid)
	rm .godoc.pid

open-docs:
	open http://localhost:6064/pkg/github.com/$(REPO_OWNER)/migrations

# example: make release V=0.0.0
release:
	git tag v$(V)
	@read -p "Press enter to confirm and push to origin ..." && git push origin v$(V)


.PHONY: build-cli clean test-short test test-with-flags deps html-coverage \
        restore-import-paths rewrite-import-paths list-external-deps release \
        docs kill-docs open-docs kill-orphaned-docker-containers dep-ensure \
        dep-update

SHELL = /bin/bash
RAND = $(shell echo $$RANDOM)
