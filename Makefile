export GOBIN := ${CURDIR}/bin
GO_PROJ      := github.com/unstoppablemango/lang

DOTNET_BUILD_CONFIG := Debug

LOCALBIN := ${CURDIR}/bin
DEVCTL   := ${LOCALBIN}/devctl
DPRINT   := ${LOCALBIN}/dprint
BUF      := ${LOCALBIN}/buf
DOTNET   := ${LOCALBIN}/dotnet
FANTOMAS := ${LOCALBIN}/fantomas

build: .make/dotnet-build bin/ir
gen: .make/buf-gen
test: .make/dotnet-test .make/ginkgo-test
format: .make/fantomas-format .make/dprint-format .make/buf-format
tidy: go.sum
dev: .envrc bin/devctl bin/dotnet

go.sum: go.mod $(shell $(DEVCTL) list --go)
	go mod tidy

go.mod:
	go mod init ${GO_PROJ}

.config/dotnet-tools.json:
	$(DOTNET) new tool-manifest

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

bin/dotnet: | .make/dotnet
	rm -f $@ && ln -s ${CURDIR}/.make/dotnet/dotnet $@

bin/fantomas: .config/dotnet-tools.json
	dotnet tool restore
	printf '#!/bin/bash\ndotnet fantomas $$@\n' > $@ && chmod +x $@

bin/dprint: .versions/dprint | .make/dprint/install.sh bin
	DPRINT_INSTALL=${CURDIR} .make/dprint/install.sh $(shell $(DEVCTL) v dprint)
	@touch $@

bin/buf: .versions/buf | bin/devctl
	go install github.com/bufbuild/buf/cmd/buf@$(shell $(DEVCTL) $<)

src/UnMango.Lang.Host/bin/lang-host: $(shell $(DEVCTL) list --cs) | bin/devctl
	dotnet publish src/UnMango.Lang.Host -p:DebugSymbols=false \
	--use-current-runtime --self-contained --configuration ${DOTNET_BUILD_CONFIG} --output $(dir $@)

.make/dotnet-install.sh: | .make
	curl -fsSL https://dot.net/v1/dotnet-install.sh > $@ && chmod +x $@

.make/dotnet: global.json | .make/dotnet-install.sh
	.make/dotnet-install.sh --install-dir $@ --jsonfile $< --no-path

.make/dotnet-build: $(shell $(DEVCTL) list --dotnet) Lang.sln | .make bin/devctl
	$(DOTNET) build
	@touch $@

.make/dotnet-test: $(shell $(DEVCTL) list --dotnet) Lang.sln | .make bin/devctl
	$(DOTNET) test
	@touch $@

.make/dotnet-format: $(shell $(DEVCTL) list --cs) | .make bin/devctl
	$(DOTNET) format --include $?

.make/fantomas-format: $(shell $(DEVCTL) list --fs) | .make bin/fantomas
	$(FANTOMAS) $?
	@touch $@

.make/ginkgo-test: $(shell $(DEVCTL) list --go) | .make bin/devctl bin/ginkgo bin/lang-host
	ginkgo run -r ./
	@touch $@

.make/dprint/install.sh: | .make
	mkdir -p $(dir $@)
	curl -fsSL https://dprint.dev/install.sh -o $@
	chmod +x $@

.make/dprint-format: README.md .dprint.jsonc .github/renovate.json | bin/dprint
	$(DPRINT) fmt
	@touch $@

.make/buf-gen: $(shell $(DEVCTL) list --proto) | .make bin/buf bin/devctl
	$(BUF) generate

.make/buf-format: $(shell $(DEVCTL) list --proto) | .make bin/buf bin/devctl
	$(BUF) format --write
	@touch $@
