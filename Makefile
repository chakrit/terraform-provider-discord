#!/usr/bin/make

VERSION := v0.0.4-p9fix0
BINNAME := terraform-provider-discord_$(VERSION)
ARCH    := darwin_arm64
KEYNAME := terraform

default: bin zip shasum sig
bin: $(BINNAME)
zip: $(BINNAME)_$(ARCH).zip
shasum: $(BINNAME)_SHA256SUMS
sig: $(BINNAME)_SHA256SUMS.sig

$(BINNAME):
	go build -o $@
$(BINNAME)_$(ARCH).zip:
	zip "$@" $(BINNAME)
$(BINNAME)_SHA256SUMS:
	shasum -a 256 *.zip > "$@"
$(BINNAME)_SHA256SUMS.sig: $(BINNAME)_SHA256SUMS
	gpg --default-key "$(KEYNAME)" --detach-sign "$<"

