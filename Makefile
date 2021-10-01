#!/usr/bin/make

VERSION := 0.0.5
ARCH    := darwin_arm64
NAME    := discord
KEYNAME := terraform

SRCFILES := main.go $(wildcard **/*.go)

BINNAME  := terraform-provider-$(NAME)_v$(VERSION)
ZIPNAME  := terraform-provider-$(NAME)_$(VERSION)_$(ARCH).zip
SUMNAME  := terraform-provider-$(NAME)_$(VERSION)_SHA256SUMS
SIGNAME  := terraform-provider-$(NAME)_$(VERSION)_SHA256SUMS.sig

default: bin zip shasum sig
bin: $(BINNAME)
zip: $(ZIPNAME)
shasum: $(SUMNAME)
sig: $(SIGNAME)

$(BINNAME): $(SRCFILES)
	go build -o $@
$(ZIPNAME): $(BINNAME)
	zip "$@" "$<"
$(SUMNAME): $(ZIPNAME)
	shasum -a 256 "$<" > "$@"
$(SIGNAME): $(SUMNAME)
	gpg --default-key "$(KEYNAME)" --detach-sign "$<"

