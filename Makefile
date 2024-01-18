
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
	fyne package --os $(GOOS) --name setup --appID in.kondrash.examen --appVersion 0.0.1 --icon ../../resources/examen.png --release --sourceDir ./cmd/setup

cmd/setup/embed/install$(EXE).gz: cmd/install/install$(EXE)
	gzip -fc cmd/install/install.exe > cmd/setup/embed/install.exe.gz

cmd/setup/embed/opengl32.dll.gz: resources/opengl32.dll
	gzip -fc resources/opengl32.dll  > cmd/setup/embed/opengl32.dll.gz

cmd/install/install.exe: cmd/install/embed/opengl32.dll.gz cmd/install/embed/examen.exe.gz cmd/install/embed/examensvc.exe.gz $(wildcard cmd/install/*.go)
	fyne package --os $(GOOS) --name install --appID in.kondrash.examen --appVersion 0.0.1 --icon ../../resources/examen.png --release --sourceDir ./cmd/install

cmd/install/embed/opengl32.dll.gz: resources/opengl32.dll
	gzip -fc resources/opengl32.dll  > cmd/install/embed/opengl32.dll.gz

cmd/install/embed/examensvc.exe.gz: cmd/examensvc/examensvc.exe
	gzip -fc cmd/examensvc/examensvc.exe  > cmd/install/embed/examensvc.exe.gz

cmd/install/embed/examen.exe.gz: cmd/examen/examen.exe
	gzip -fc cmd/examen/examen.exe  > cmd/install/embed/examen.exe.gz

cmd/examen/examen.exe: $(wildcard cmd/examen/*.go)
	fyne package --os $(GOOS) --name examen --appID in.kondrash.examen --appVersion 0.0.1 --icon ../../resources/examen.png --release --sourceDir ./cmd/examen

cmd/examensvc/examensvc.exe: $(wildcard cmd/examensvc/*.go)
	fyne package --os $(GOOS) --name examensvc --appID in.kondrash.examen --appVersion 0.0.1 --icon ../../resources/examen.png --release --sourceDir ./cmd/examensvc

clean: tidy
	rm -f setup.zip

tidy:
	rm -f cmd/examen/examen.exe cmd/examensvc/examensvc.exe cmd/install/install.exe cmd/setup/setup.exe setup.zip
