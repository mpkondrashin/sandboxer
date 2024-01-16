
.PHONY: clean tidy

#ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
GOOS="unknown_OS"
MOVE="unknown_move_command"
EXE=""
TAR=""
ZIP=""

ifeq ($(OS),Windows_NT)
	GOOS=windows
	MOVE=move /Y
	EXE=.exe
	TAR="rem"
	ZIP="powershell Compress-Archive

# $(1) - archive name
# $(2) - file to put into archive`
define zip
	powershell Compress-Archive  -Force "$(2)" "$(1)"
endef
else
	GOOS=darwin
	MOVE=mv -f
	EXE=.app
	TAR="tag cfv"
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
	fyne package --os $(GOOS) --name setup --appID in.kondrash.examen --appVersion 0.0.1 --icon ../../resources/examen.png --release --sourceDir ./cmd/setup

cmd/setup/embed/install$(EXE).gz: cmd/install/install$(EXE)
	gzip -fc cmd/install/install.exe > cmd/setup/embed/install.exe.gz

cmd/setup/embed/opengl32.dll.gz: resources/opengl32.dll
	gzip -fc resources/opengl32.dll  > cmd/setup/embed/opengl32.dll.gz

cmd/install/install.exe: examen.exe examensvc.exe $(wildcard cmd/install/*.go)
	fyne package --os $(GOOS) --name install --appID in.kondrash.examen --appVersion 0.0.1 --icon ../../resources/examen.png --release --sourceDir ./cmd/install

examen.exe: $(wildcard cmd/examen/*.go)
	go build ./cmd/examen

examensvc.exe: $(wildcard cmd/examensvc/*.go)
	go build ./cmd/examensvc

clean: tidy
	rm setup.exe

tidy:
	rm examen.exe examensvc.exe install.exe 
