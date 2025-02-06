export GOBIN := ${CURDIR}/bin
GO_PROJ      := github.com/unstoppablemango/lang

DEVCTL := ${CURDIR}/bin/devctl

build: .make/dotnet-build
test: .make/dotnet-test
tidy: go.sum
dev: .envrc bin/devctl

go.sum: go.mod | bin/devctl
	go mod tidy

go.mod:
	go mod init ${GO_PROJ}

.envrc: hack/example.envrc
	cp $< $@

.make .versions bin:
	mkdir -p $@

bin/devctl: .versions/devctl | bin
	go install github.com/unmango/devctl/cmd@v$(shell cat $<)
	mv bin/cmd $@

.make/dotnet-build: $(shell $(DEVCTL) list --dotnet) Lang.sln | .make bin/devctl
	dotnet build
	@touch $@

.make/dotnet-test: $(shell $(DEVCTL) list --dotnet) Lang.sln | .make bin/devctl
	dotnet test
	@touch $@
