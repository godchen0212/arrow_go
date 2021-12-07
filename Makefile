# Copyright (C) 2019-2020 Zilliz. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance
# with the License. You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software distributed under the License
# is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
# or implied. See the License for the specific language governing permissions and limitations under the License.

GO		?= go
PWD 	:= $(shell pwd)
GOPATH 	:= $(shell $(GO) env GOPATH)

INSTALL_PATH := $(PWD)/bin
LIBRARY_PATH := $(PWD)/lib

all: build-cpp build-go


fmt:
ifdef GO_DIFF_FILES
	@echo "Running $@ check"
	@GO111MODULE=on env bash $(PWD)/scripts/gofmt.sh $(GO_DIFF_FILES)
else
	@echo "Running $@ check"
	@GO111MODULE=on env bash $(PWD)/scripts/gofmt.sh cmd/
	@GO111MODULE=on env bash $(PWD)/scripts/gofmt.sh internal/
	@GO111MODULE=on env bash $(PWD)/scripts/gofmt.sh tests/go/
endif


# Build various components locally.
binlog:
	@echo "Building binlog ..."
	@mkdir -p $(INSTALL_PATH) && go env -w CGO_ENABLED="1" && GO111MODULE=on $(GO) build -o $(INSTALL_PATH)/binlog $(PWD)/cmd/tools/binlog/main.go 1>/dev/null

BUILD_TAGS = $(shell git describe --tags --always --dirty="-dev")
BUILD_TIME = $(shell date --utc)
GIT_COMMIT = $(shell git rev-parse --short HEAD)
GO_VERSION = $(shell go version)

my_example: build-cpp
	@echo "Building my_example ..."
	@mkdir -p $(INSTALL_PATH) && go env -w CGO_ENABLED="1" && GO111MODULE=on $(GO) build \
		-ldflags="-X 'main.BuildTags=$(BUILD_TAGS)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.GoVersion=$(GO_VERSION)'" \
		-o $(INSTALL_PATH)/my_example $(PWD)/cmd/main.go 1>/dev/null

build-go: my_example

build-cpp:
	@echo "Building my_example cpp library ..."
	@(env bash $(PWD)/scripts/cwrapper_build.sh -t Debug -f "$(CUSTOM_THIRDPARTY_PATH)")


# Build each component and install binary to $GOPATH/bin.
install: all
	@echo "Installing binary to './bin'"
	@mkdir -p $(GOPATH)/bin && cp -f $(PWD)/bin/my_example $(GOPATH)/bin/my_example
#	@mkdir -p $(LIBRARY_PATH) && cp -P $(PWD)/internal/core/output/lib/* $(LIBRARY_PATH)
	@echo "Installation successful."

clean:
	@echo "Cleaning up all the generated files"
	@find . -name '*.test' | xargs rm -fv
	@find . -name '*~' | xargs rm -fv
	@rm -rf bin/
	@rm -rf lib/
	@rm -rf $(GOPATH)/bin/*
