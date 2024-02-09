
.PHONY: clean tidy

#GOOPTS := -ldflags="-extldflags=-static"
# -tags sqlite_omit_load_extension
GOOS=windows

REV=$(shell git rev-list --tags --max-count=1)
VERSION_FULL=$(shell git describe --tags $(REV))
VERSION_V := $(word 1,$(subst -, ,$(VERSION_FULL)))
VERSION := $(subst v,,$(VERSION_V))
ifeq ($(VERSION),)
VERSION := v
endif
BUILD = $(shell git rev-list --all --count)
ifeq ($(BUILD),)
BUILD := b
endif

EXE=.exe

define zip
	powershell Compress-Archive  -Force "$(2)" "$(1)"
endef

setup.zip: cmd/setup/setup.exe
	$(call zip, "setup.zip" , "cmd/setup/setup.exe")

preproc.exe:  $(wildcard cmd/preproc/*.go)
	go build ./cmd/preproc

cmd/setup/setup.exe.manifest: cmd/preproc/preproc.exe cmd/setup/manifest.template
	preproc.exe --version $(VERSION) --build $(BUILD) $< $@

cmd/setup/setup.syso: cmd/setup/setup.exe.manifest
	rsrc -manifest cmd/setup/setup.exe.manifest -o cmd/setup/setup.syso

cmd/setup/setup.exe: cmd/setup/embed/install.exe.gz cmd/setup/embed/opengl32.dll.gz $(wildcard cmd/setup/*.go)
	fyne package --os $(GOOS) --name setup --appID in.kondrash.sandboxer --appVersion 0.0.1 --icon ../../resources/icon.png --release --sourceDir ./cmd/setup

cmd/setup/embed/install.exe.gz: cmd/install/install.exe
	gzip -fc cmd/install/install.exe > cmd/setup/embed/install.exe.gz

cmd/setup/embed/opengl32.dll.gz: resources/opengl32.dll
	gzip -fc resources/opengl32.dll  > cmd/setup/embed/opengl32.dll.gz

cmd/install/install.exe: cmd/install/embed/opengl32.dll.gz cmd/install/embed/sandboxer.exe.gz cmd/install/embed/submit.exe.gz $(wildcard cmd/install/*.go) $(wildcard pkg/extract/*.go)  $(wildcard pkg/globals/*.go) cmd/install/resource.go
	fyne package --os $(GOOS) --name install --appID in.kondrash.sandboxer --appVersion 0.0.1 --icon ../../resources/icon.png --release --sourceDir ./cmd/install

cmd/install/resource.go: resources/icon_transparent.png 
	fyne bundle --name ApplicationIcon --package main --output cmd/install/resource.go resources/icon_transparent.png 

cmd/install/embed/opengl32.dll.gz: resources/opengl32.dll
	gzip -fc resources/opengl32.dll  > cmd/install/embed/opengl32.dll.gz

cmd/install/embed/sandboxer.exe.gz: cmd/sandboxer/sandboxer.exe
	gzip -fc cmd/sandboxer/sandboxer.exe  > cmd/install/embed/sandboxer.exe.gz

cmd/submit/submit.exe: $(wildcard cmd/submit/*.go)  $(wildcard pkg/globals/*.go) 
	fyne package --os $(GOOS) --name submit --appID in.kondrash.sandboxer --appVersion 0.0.1 --icon ../../resources/icon.png --release --sourceDir ./cmd/submit

cmd/install/embed/submit.exe.gz: cmd/submit/submit.exe
	gzip -fc cmd/submit/submit.exe  > cmd/install/embed/submit.exe.gz

cmd/sandboxer/sandboxer.exe: $(wildcard cmd/sandboxer/*.go) $(wildcard pkg/globals/*.go) cmd/sandboxer/icon.go
	fyne package --os $(GOOS) --name sandboxer --appID in.kondrash.sandboxer --appVersion 0.0.1 --icon ../../resources/icon.png --release --sourceDir ./cmd/sandboxer

cmd/sandboxer/icon.go: resources/icon.png 
	fyne bundle --name ApplicationIcon --package main --output cmd/sandboxer/icon.go resources/icon.png 

clean: tidy
	rm -f setup.zip

tidy:
	rm -f cmd/sandboxer/sandboxer.exe cmd/install/install.exe cmd/setup/setup.exe

