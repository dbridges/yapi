.PHONY: all run run-with-race clean install uninstall

VERSION=0.0.2
GO_SRC = $(shell find . -iname '*.go' ! -iname "*test.go")
BINDIR?=/usr/local/bin
BINNAME?=yapi
OUTPUT_PATH=$(BINNAME)

all: $(OUTPUT_PATH)

$(OUTPUT_PATH): $(GO_SRC)
	go build -ldflags "-X github.com/dbridges/yapi/cmd.VersionString=$(VERSION)" -o $@

run:
	@go run $(GO_SRC)

run-with-race:
	@GORACE="log_path=race_log" go run -race *.go

clean:
	rm yapi
	rm -f race_log.*

install: $(OUTPUT_PATH)
	mkdir -p $(BINDIR)
	install $(OUTPUT_PATH) $(BINDIR)/$(BINNAME)

uninstall:
	rm -f $(BINDIR)/$(BINNAME)
