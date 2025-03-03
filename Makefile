export GOBIN := ${CURDIR}/bin
GO_PROJ      := github.com/unstoppablemango/lang

DOTNET_CONFIG := Debug

LOCALBIN := ${CURDIR}/bin
DEVCTL   ?= go tool devctl
DPRINT   ?= ${LOCALBIN}/dprint
BUF      ?= ${LOCALBIN}/buf
DOTNET   ?= ${LOCALBIN}/dotnet
FANTOMAS ?= ${LOCALBIN}/fantomas
GINKGO   ?= go tool ginkgo
NINJA    ?= ${LOCALBIN}/ninja

ifeq (${CI},)
TEST_FLAGS := --label-filter !E2E
else
TEST_FLAGS := --github-output --trace --cover
endif

all: bin/lang-host bin/ir
gen: .make/buf-gen
test: .make/dotnet-test .make/ginkgo-test
format: .make/fantomas-format .make/dotnet-format .make/dprint-format .make/buf-format
tidy: go.sum
dev: .envrc bin/devctl bin/dotnet
ci: .make test build build/lang

clean: .make/dotnet-clean
	rm -rf src/**/{bin,obj}

.PHONY: build build/lang
build: | bin/ninja
	cmake --preset default -DCMAKE_MAKE_PROGRAM=$(NINJA)

build/lang:
	cmake --build build

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

bin/dotnet: | .make/dotnet
	rm -f $@ && ln -s ${CURDIR}/.make/dotnet/dotnet $@

bin/fantomas: .config/dotnet-tools.json
	dotnet tool restore
	printf '#!/bin/bash\ndotnet fantomas $$@\n' > $@ && chmod +x $@

bin/dprint: .versions/dprint | .make/dprint/install.sh bin
	DPRINT_INSTALL=${CURDIR} .make/dprint/install.sh $(shell $(DEVCTL) v dprint)
	@touch $@

bin/buf: .versions/buf
	go install github.com/bufbuild/buf/cmd/buf@$(shell $(DEVCTL) $<)

bin/ninja: | .make/ninja.zip
	unzip ${CURDIR}/$| -d ${LOCALBIN}

bin/vcpkg: | tools/vcpkg/vcpkg
	ln -s ${CURDIR}/$| ${CURDIR}/$@

src/UnMango.Lang.Host/bin/lang-host: $(shell $(DEVCTL) list --cs)
	dotnet publish src/UnMango.Lang.Host -p:DebugSymbols=false \
	--use-current-runtime --self-contained --configuration ${DOTNET_CONFIG} --output $(dir $@)

tools/vcpkg/vcpkg: tools/vcpkg/bootstrap-vcpkg.sh
	$< --disableMetrics

tools/vcpkg/bootstrap-vcpkg.sh:
	git submodule update --init --recursive

.make/dotnet-install.sh: | .make
	curl -fsSL https://dot.net/v1/dotnet-install.sh > $@ && chmod +x $@

.make/dotnet: global.json | .make/dotnet-install.sh
	.make/dotnet-install.sh --install-dir $@ --jsonfile $< --no-path

.make/dotnet-build: $(shell $(DEVCTL) list --dotnet) Lang.sln | .make bin/dotnet
	$(DOTNET) build
	@touch $@

.make/dotnet-test: $(shell $(DEVCTL) list --dotnet) Lang.sln | .make bin/dotnet
	$(DOTNET) test
	@touch $@

.make/dotnet-format: $(shell $(DEVCTL) list --cs) | .make bin/dotnet
	$(DOTNET) format --include $?
	@touch $@

.make/dotnet-clean: | .make bin/dotnet
	$(DOTNET) clean
	@touch $@

.make/fantomas-format: $(shell $(DEVCTL) list --fs) | .make bin/fantomas
	$(FANTOMAS) $?
	@touch $@

.make/ginkgo-test: $(shell $(DEVCTL) list --go) | .make bin/lang-host
	$(GINKGO) run ${TEST_FLAGS} $(sort $(dir $?))
	@touch $@

.make/dprint/install.sh: | .make
	mkdir -p $(dir $@)
	curl -fsSL https://dprint.dev/install.sh -o $@
	chmod +x $@

.make/dprint-format: README.md .dprint.jsonc .github/renovate.json | bin/dprint
	$(DPRINT) fmt $?
	@touch $@

.make/buf-gen: $(shell $(DEVCTL) list --proto) | .make bin/buf bin/devctl
	$(BUF) generate

.make/buf-format: $(shell $(DEVCTL) list --proto) | .make bin/buf bin/devctl
	$(BUF) format --write
	@touch $@

.make/ninja.zip: .versions/ninja
	curl -Lo $@ https://github.com/ninja-build/ninja/releases/download/$(shell $(DEVCTL) $<)/ninja-$(shell go env GOOS).zip
