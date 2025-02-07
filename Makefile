export GOBIN := ${CURDIR}/bin
GO_PROJ      := github.com/unstoppablemango/lang

LOCALBIN := ${CURDIR}/bin
DEVCTL   := ${LOCALBIN}/devctl
DPRINT   := ${LOCALBIN}/dprint
BUF      := ${LOCALBIN}/buf

build: .make/dotnet-build bin/ir
gen: .make/buf-gen
test: .make/dotnet-test
format: .make/dprint-format
tidy: go.sum
dev: .envrc bin/devctl

go.sum: go.mod $(shell $(DEVCTL) list --go)
	go mod tidy

go.mod:
	go mod init ${GO_PROJ}

.envrc: hack/example.envrc
	cp $< $@

.make bin:
	mkdir -p $@

bin/ir: $(shell $(DEVCTL) list --go)
	go build -o $@ -tags=llvm19 ./cmd/ir

bin/devctl: .versions/devctl | bin
	go install github.com/unmango/devctl/cmd@v$(shell cat $<)
	mv bin/cmd $@

bin/dprint: .versions/dprint | .make/dprint/install.sh bin
	DPRINT_INSTALL=${CURDIR} .make/dprint/install.sh $(shell $(DEVCTL) v dprint)
	@touch $@

bin/buf: .versions/buf | bin
	go install github.com/bufbuild/buf/cmd/buf@$(shell $(DEVCTL) $<)

.make/dotnet-build: $(shell $(DEVCTL) list --dotnet) Lang.sln | .make bin/devctl
	dotnet build
	@touch $@

.make/dotnet-test: $(shell $(DEVCTL) list --dotnet) Lang.sln | .make bin/devctl
	dotnet test
	@touch $@

.make/dprint/install.sh:
	mkdir -p $(dir $@)
	curl -fsSL https://dprint.dev/install.sh -o $@
	chmod +x $@

.make/dprint-format: README.md .dprint.jsonc .github/renovate.json | bin/dprint
	$(DPRINT) fmt

.make/buf-gen:
	$(BUF) generate
