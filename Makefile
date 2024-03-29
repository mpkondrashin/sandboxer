#
# Sandboxer (c) 2024 by Mikhail Kondrashin (mkondrashin@gmail.com)
# Software is distributed under MIT license as stated in LICENSE file
#
# Makefile
#
# Makefile for Windows
#

.PHONY: clean tidy

#GOOPTS := -ldflags="-extldflags=-static"
GOOPTS=-race
# -tags sqlite_omit_load_extension
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

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
#BUILD_OPTS :=  -ldflags="-X 'sandboxer/pkg/globals.Version=v1.1.1'"
#-ldflags "-X 'github.com/mpkondrashin/sandboxer/pkg/globals.Version=$(VERSION)'"
# "-X 'github.com/mpkondrashin/sandboxer/pkg/globals.Build=$(BUILD)'"

define zip
	powershell Compress-Archive  -Force "$(2)" "$(1)"
endef

setup_$(GOOS)_$(GOARCH).zip: cmd/setup/setup.exe
	$(call zip, setup_$(GOOS)_$(GOARCH).zip , "cmd/setup/setup.exe")

cmd/setup/setup.exe: cmd/setup/embed/install.exe.gz cmd/setup/embed/opengl32.dll.gz  $(wildcard cmd/setup/*.go) $(wildcard pkg/*/*.go) pkg/globals/version.go cmd/setup/setup.syso
	GOOS=windows go build -C ./cmd/setup -ldflags -H=windowsgui 

cmd/setup/embed/opengl32.dll.gz: resources/opengl32.dll
	gzip -fc $< > $@

cmd/setup/setup.exe.manifest: preproc.exe cmd/setup/manifest.template
	GOOS=windows GOARCH=amd64 ./preproc.exe --version $(VERSION) --build $(BUILD) cmd/setup/manifest.template cmd/setup/setup.exe.manifest

preproc.exe:  $(wildcard cmd/preproc/*.go)
	go build ./cmd/preproc

cmd/setup/setup.syso: cmd/setup/setup.exe.manifest
	go get -u github.com/akavel/rsrc
	go install github.com/akavel/rsrc
	rsrc -manifest ./cmd/setup/setup.exe.manifest -o ./cmd/setup/setup.syso

#--icon ../../resources/icon.png

cmd/setup/embed/install.exe.gz: cmd/install/install.exe
	gzip -fc $< > $@

cmd/install/install.exe: cmd/install/embed/sandboxer.tar.gz cmd/install/embed/LICENSE $(wildcard cmd/install/*.go) $(wildcard pkg/*/*.go) pkg/globals/version.go resources/icon.png cmd/install/resource.go
	fyne package --os $(GOOS) --name install --appID in.kondrash.sandboxer --appVersion $(VERSION) --appBuild $(BUILD) --icon ../../resources/icon.png --release --sourceDir ./cmd/install

cmd/install/embed/LICENSE: LICENSE
	cp $< $@

cmd/install/resource.go: resources/icon_transparent.png 
	fyne bundle --name ApplicationIcon --package main --output cmd/install/resource.go resources/icon_transparent.png 

cmd/install/embed/sandboxer.tar.gz: cmd/submit/submit.exe cmd/sandboxer/sandboxer.exe LICENSE resources/opengl32.dll resources/DroidSansHebrew-Regular.ttf resources/icon_transparent.png
	tar cfvz $@ LICENSE -C resources opengl32.dll DroidSansHebrew-Regular.ttf icon_transparent.png -C ../cmd/submit submit.exe -C ../../cmd/sandboxer sandboxer.exe 

cmd/submit/submit.exe: $(wildcard cmd/submit/*.go)  $(wildcard pkg/*/*.go) pkg/globals/version.go resources/icon.png
	fyne package --os $(GOOS) --name submit --appID in.kondrash.sandboxer --appVersion $(VERSION) --appBuild $(BUILD) --icon ../../resources/icon.png --release --sourceDir ./cmd/submit

cmd/sandboxer/sandboxer.exe: $(wildcard cmd/sandboxer/*.go) $(wildcard pkg/*/*.go) pkg/globals/version.go cmd/sandboxer/icon.go
	fyne package --os $(GOOS) --name sandboxer --appID in.kondrash.sandboxer --appVersion $(VERSION) --appBuild $(BUILD) --icon ../../resources/icon.png --release --sourceDir ./cmd/sandboxer

cmd/sandboxer/icon.go: resources/icon.png 
	fyne bundle --name ApplicationIcon --package main --output cmd/sandboxer/icon.go resources/icon.png 

pkg/globals/version.go: cmd/genver/main.go
	go run ./cmd/genver/main.go $(VERSION_FULL) $(BUILD) pkg/globals/version.go

clean: cleansetup celaninstall
	rm -f setup.zip cmd/sandboxer/sandboxer.exe cmd/submit/submit.exe preproc.exe pkg/globals/version.go

cleansetup:
	rm -f cmd/setup/setup.exe cmd/setup/setup.exe.manifest cmd/setup/setup.syso cmd/setup/sandboxer_setup_wizard.log cmd/setup/sandboxer_setup.log cmd/setup/install.exe cmd/setup/embed/install.exe.gz cmd/setup/embed/opengl32.dll.gz 

celaninstall:
	rm -f cmd/install/install.exe cmd/install/embed/sandboxer.tar.gz