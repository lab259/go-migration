VERSION ?= $(shell git describe --tags 2>/dev/null | cut -c 2-)
TEST_FLAGS ?=

COVERDIR=$(CURDIR)/.cover
COVERAGEFILE=$(COVERDIR)/cover.out

test:
	@ginkgo --failFast ./...

test-watch:
	@ginkgo watch -cover -r ./...

coverage:
	@mkdir -p $(COVERDIR)
	@ginkgo -r -covermode=count --cover --trace ./
	@echo "mode: count" > "${COVERAGEFILE}"
	@find . -type f -name *.coverprofile -exec grep -h -v "^mode:" {} >> "${COVERAGEFILE}" \; -exec rm -f {} \;

coverage-ci:
	@mkdir -p $(COVERDIR)
	@ginkgo -r -covermode=count --cover --trace ./
	@echo "mode: count" > "${COVERAGEFILE}"
	@find . -type f -name *.coverprofile -exec grep -h -v "^mode:" {} >> "${COVERAGEFILE}" \; -exec rm -f {} \;

coverage-html:
	@go tool cover -html="${COVERAGEFILE}" -o .cover/report.html

dcup:
	@docker-compose up -d

dcdn:
	@docker-compose down --remove-orphans

vet:
	@go vet ./...

lint:
	@golint

fmt:
	@go fmt ./...

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
	open http://localhost:6064/pkg/github.com/lab259/go-migration

# example: make release V=0.0.0
release:
	git tag v$(V)
	@read -p "Press enter to confirm and push to origin ..." && git push origin v$(V)


.PHONY: test test-watch coverage coverage-ci coverage-html dcup dcdn vet lint fmt docs kill-docs open-docs release