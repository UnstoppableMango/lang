export GOBIN := ${CURDIR}/bin
GO_PROJ      := github.com/unstoppablemango/lang

DOTNET_RID ?= $(shell dotnet --info | grep RID | sed 's/RID\:\s*//g')

LOCALBIN := ${CURDIR}/bin
DEVCTL   := ${LOCALBIN}/devctl
DPRINT   := ${LOCALBIN}/dprint
BUF      := ${LOCALBIN}/buf

build: .make/dotnet-build bin/ir
gen: .make/buf-gen
test: .make/dotnet-test .make/ginkgo-test
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

bin/lang-host: src/UnMango.Lang.Host/bin/lang-host
	cp $< $@

bin/devctl: .versions/devctl | bin
	go install github.com/unmango/devctl/cmd@v$(shell cat $<)
	mv bin/cmd $@

bin/ginkgo: go.mod | bin
	go install github.com/onsi/ginkgo/v2/ginkgo

bin/dprint: .versions/dprint | .make/dprint/install.sh bin
	DPRINT_INSTALL=${CURDIR} .make/dprint/install.sh $(shell $(DEVCTL) v dprint)
	@touch $@

bin/buf: .versions/buf | bin/devctl
	go install github.com/bufbuild/buf/cmd/buf@$(shell $(DEVCTL) $<)

src/UnMango.Lang.Host/bin/lang-host: $(shell $(DEVCTL) list --cs) | bin/devctl
	dotnet publish src/UnMango.Lang.Host -p:DebugSymbols=false \
	--use-current-runtime --self-contained --configuration Release --output $(dir $@)

.make/dotnet-build: $(shell $(DEVCTL) list --dotnet) Lang.sln | .make bin/devctl
	dotnet build
	@touch $@

.make/dotnet-test: $(shell $(DEVCTL) list --dotnet) Lang.sln | .make bin/devctl
	dotnet test
	@touch $@

.make/ginkgo-test: $(shell $(DEVCTL) list --go) | .make bin/devctl bin/ginkgo
	ginkgo run -r ./
	@touch $@

.make/dprint/install.sh: | .make
	mkdir -p $(dir $@)
	curl -fsSL https://dprint.dev/install.sh -o $@
	chmod +x $@

.make/dprint-format: README.md .dprint.jsonc .github/renovate.json | bin/dprint
	$(DPRINT) fmt

.make/buf-gen: | .make
	$(BUF) generate
