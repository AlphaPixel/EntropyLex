# MIT License
#
# Copyright 2026 arcadium.dev <info@arcadium.dev>
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

export SHELL := /bin/bash

ifeq ($(shell uname -m),arm64)
  arch := arm64
else
  arch := amd64
endif

os := $(shell uname -s | tr [A-Z] [a-z])

# sha_len is the length of the sha sum used with the version and the sha
sha_len := 7

# version is the version of the current branch. For code that matches a
# released version we want the exact version match, i.e. v1.0.0. For code that
# is part of work in progress we want a version that denotes a path to a
# release version, i.e. v1.0.0-5-g07a65db-dirty, where the closest release is
# v1.0.0, the 5 denotes that the code is 5 commits ahead of the release,
# g07a65db is the git sha of the latest commit, and dirty denotes that there
# are uncommitted changes to the code.
export version := $(shell basename $(shell git describe --tags --dirty --abbrev=$(sha_len)))

# branch is the name of the current branch
export branch ?= $(shell git rev-parse --abbrev-ref HEAD)

# commit is the shasum of the latest commit
export commit := $(shell git rev-parse HEAD)

# date is the date of the build
export date := $(shell date -u -Iseconds)

# ldflags are the go linker flags we pass to the go command.
#   -s    Omit the symbol table and debug information.
#   -w    Omit the DWARF symbol table.
#   -X importpath.name=value
#         Set the value of the string variable in import path named name to
#         value.  This is only effective if the variable is declared in the
#         source code either uninitialized or initialized to a constant string
#         expression.
export ldflags := -s -w \
	-X main.Version=$(version) \
	-X main.Branch=$(branch) \
	-X main.Commit=$(commit) \
	-X main.Date=$(date)


entropylex := dist/$(os)/$(arch)/entropylex
cmds := $(entropylex)

schema := entropylex/dictionary/schema/lxj-v1.schema.json

el8dict := entropylex/cmd/entropylex/dict/entropylex-en-8-test-v1.lxj
dicts := $(el8dict)

# ____ all __________________________________________________________________

.PHONY: all

all: $(dicts) test lint $(cmds)


$(schema): ../../../data/dict/lxj-v1.schema.json
	@printf "\nCopying lxj-v1.schema.json locally...\n"
	@cp $< $@

$(el8dict): ../../../tests/fixtures/dict/entropylex-en-8-test-v1.lxj
	@printf "\nCopying entropylex-en-8-test-v1.lxj locally...\n"
	cp $< $@

# ____ lint __________________________________________________________________

.PHONY: fmt tidy vet staticcheck vuln lint

fmt:
	@printf "\nRunning go fmt...\n"
	go fmt ./...

tidy:
	@printf "\nRunning go mod tidy...\n"
	go mod tidy

vet:
	@printf "\nRunning go vet...\n"
	go vet ./...

staticcheck:
	@printf "\nRunning staticcheck...\n"
	go tool staticcheck ./...

vuln:
	@printf "\nRunning govulncheck...\n"
	go tool govulncheck ./...

lint: fmt tidy vet staticcheck vuln
	@printf "\nChecking for changed files...\n"
	git status --porcelain
	@printf "\n"
	@if [[ "$${CI}" == "true" ]]; then $$(exit $$(git status --porcelain | wc -l)); fi

# ____ test __________________________________________________________________

.PHONY: unit_test test

unit_test:
	@printf "\nRunning go test...\n"
	go test -count=1 -cover -race ./...

test: unit_test

# ____ cmds ___________________________________________________________________

.PHONY: $(entropylex)

$(entropylex): dist/$(os)/$(arch)
	GOOS="$(os)" GOARCH="$(arch)" CGO_ENAGLED=0 \
	  go build -ldflags="$(ldflags)" -o $@ ./entropylex/cmd/entropylex

dist/$(os)/$(arch):
	-mkdir -p $@

# ____ clean artifacts _______________________________________________________

.PHONY: clean

clean:
	@printf "\nClean...\n"
	-go clean -i -r -cache -testcache -modcache -fuzzcache
	-rm -rf dist
