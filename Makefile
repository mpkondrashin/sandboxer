
.PHONY: clean tidy

#ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)
OS="unknown_OS"
MOVE="unknown_move_command"
EXE=""
TAR=""
ZIP=""

ifeq ($(OS),Windows_NT)
	OS=windows
	MOVE=move /Y
	EXE=.exe
	TAR="rem"
	ZIP="powershell Compress-Archive

define zip
	# $(1) - archive name
	# $(2) - file to put into archive`
	#    $(eval $@_HOSTNAME = $(1))
	#    $(eval $@_PORT = $(2))
	powershell Compress-Archive "$(2)" "$(1)"
endef
else
	OS=darwin
	MOVE=mv -f
	EXE=.app
	TAR="tag cfv"
define zip
	# $(1) - archive name
	# $(2) - file to put into archive`
	#    $(eval $@_HOSTNAME = $(1))
	#    $(eval $@_PORT = $(2))
	zip $(1) $(2)
endef
endif



cmd/setup/setup$(EXE): cmd/setup/embed/install.exe.gz cmd/setup/embed/opengl32.dll.gz $(wildcard cmd/setup/*.go)
	echo "OS = $(OS)" 
	fyne package --os $(OS) --name setup --appID in.kondrash.examen --appVersion 0.0.1 --icon ../../resources/examen.png --release --sourceDir ./cmd/setup
	$(call zip, "setup.zip" , "cmd/setup/setup$(EXE)")
	#$(MOVE) cmd/setup/setup.exe .

cmd/setup/embed/install$(EXE).gz: cmd/install/install$(EXE)
	tac cfvgzip -fc cmd/install/install.exe > cmd/setup/embed/install.exe.gz
	gzip -fc cmd/install/install.exe > cmd/setup/embed/install.exe.gz

cmd/setup/embed/opengl32.dll.gz: resources/opengl32.dll
	gzip -fc resources/opengl32.dll  > cmd/setup/embed/opengl32.dll.gz

cmd/install/install.exe: examen.exe examensvc.exe $(wildcard cmd/install/*.go)
	fyne package --os $(OS) --name install --appID in.kondrash.examen --appVersion 0.0.1 --icon ../../resources/examen.png --release --sourceDir ./cmd/install

examen.exe: $(wildcard cmd/examen/*.go)
	go build ./cmd/examen

examensvc.exe: $(wildcard cmd/examensvc/*.go)
	go build ./cmd/examensvc

clean: tidy
	rm setup.exe

tidy:
	rm examen.exe examensvc.exe install.exe 
