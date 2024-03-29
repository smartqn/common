package version

/*
CURDIR := $(shell pwd)

GO        := go
GOBUILD   := GOPATH=$(GOPATH) CGO_ENABLED=0 $(GO) build $(BUILD_FLAG)
GOTEST    := GOPATH=$(GOPATH) CGO_ENABLED=1 $(GO) test -p 3

LDFLAGS += -X "github.com/smartqn/common/version.VERSION=$(shell git describe --tags --dirty)"
LDFLAGS += -X "github.com/smartqn/common/version.BUILDTIME=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
LDFLAGS += -X "github.com/smartqn/common/version.GITHASH=$(shell git rev-parse HEAD)"
LDFLAGS += -X "github.com/smartqn/common/version.GITBRANCH=$(shell git rev-parse --abbrev-ref HEAD)"

all: SmartProxy

BUILDDIR=$(CURDIR)

SmartProxy:
	@mkdir -p $(BUILDDIR)
	# $(GOBUILD) -v -o $(BUILDDIR)/$@  $(CURDIR)/main.go
	$(GOBUILD) -ldflags '$(LDFLAGS)' -o $(BUILDDIR)/$@ $(BUILDDIR)/main.go

clean:
	@rm $(BUILDDIR)/SmartProxy
*/

// sample
// {
// CURDIR := $(shell pwd)
//
// GO        := go
// GOBUILD   := GOPATH=$(GOPATH) CGO_ENABLED=0 $(GO) build $(BUILD_FLAG)
// GOTEST    := GOPATH=$(GOPATH) CGO_ENABLED=1 $(GO) test -p 3
// TARGET= SmartOutCall
//
// LDFLAGS += -X "smartqn/common/version.VERSION=$(shell git describe --tags --dirty)"
// LDFLAGS += -X "smartqn/common/version.BUILDTIME=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
// LDFLAGS += -X "smartqn/common/version.GITHASH=$(shell git rev-parse HEAD)"
// LDFLAGS += -X "smartqn/common/version.GITBRANCH=$(shell git rev-parse --abbrev-ref HEAD)"
//
// all: $(TARGET)
//
// BUILDDIR=$(CURDIR)
//
// $(TARGET):
// 	@mkdir -p $(BUILDDIR)
// 	$(GOBUILD) -ldflags '$(LDFLAGS)' -v -o $(BUILDDIR)/$@  $(CURDIR)/main.go
//
// clean:
// 	@[ -f $(BUILDDIR)/$(TARGET) ] && rm $(BUILDDIR)/$(TARGET) || true
// }
