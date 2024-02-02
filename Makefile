
.PHONY: clean tidy

#ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
GOOS="unknown_OS"
MOVE="unknown_move_command"
EXE=""
TAR=""
ZIP=""

ifeq ($(OS),Windows_NT)
	GOOS=windows
#	MOVE=move /Y
	EXE=.exe
#	TAR="rem"
	ZIP="powershell Compress-Archive

# $(1) - archive name
# $(2) - file to put into archive`
define zip
	powershell Compress-Archive  -Force "$(2)" "$(1)"
endef
else
	GOOS=darwin
#	MOVE=mv -f
	EXE=.app
#	TAR="tag cfv"
# $(1) - archive name
# $(2) - file to put into archive`
define zip
	#    $(eval $@_HOSTNAME = $(1))
	#    $(eval $@_PORT = $(2))
	zip $(1) $(2)
endef
endif

setup.zip: cmd/setup/setup$(EXE)
#	echo $(wildcard cmd/install/*.go)
	$(call zip, "setup.zip" , "cmd/setup/setup$(EXE)")

cmd/setup/setup$(EXE): cmd/setup/embed/install$(EXE).gz cmd/setup/embed/opengl32.dll.gz $(wildcard cmd/setup/*.go)
#	echo $(wildcard cmd/install/*.go)
	fyne package --os $(GOOS) --name setup --appID in.kondrash.sandboxer --appVersion 0.0.1 --icon ../../resources/icon.png --release --sourceDir ./cmd/setup

cmd/setup/embed/install$(EXE).gz: cmd/install/install$(EXE)
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
