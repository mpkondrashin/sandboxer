
.PHONY: clean tidy

#ifeq ($(PROCESSOR_ARCHITEW6432),AMD64)

ifeq ($(OS),Windows_NT)
	MOVE=move /Y
else
	MOVE=mv -f
endif


cmd/setup/setup.exe: cmd/setup/embed/install.exe.gz cmd/setup/embed/opengl32.dll.gz $(wildcard cmd/setup/*.go)
		fyne package --os windows --name setup --appID in.kondrash.examen --appVersion 0.0.1 --icon ../../resources/examen.png --release --sourceDir ./cmd/setup
		$(MOVE) cmd/setup/setup.exe .

cmd/setup/embed/install.exe.gz: cmd/install/install.exe
	gzip -fc cmd/install/install.exe > cmd/setup/embed/install.exe.gz

cmd/setup/embed/opengl32.dll.gz: resources/opengl32.dll
	gzip -fc resources/opengl32.dll  > cmd/setup/embed/opengl32.dll.gz

cmd/install/install.exe: examen.exe examensvc.exe $(wildcard cmd/install/*.go)
	fyne package --os windows --name install --appID in.kondrash.examen --appVersion 0.0.1 --icon ../../resources/examen.png --release --sourceDir ./cmd/install

examen.exe: $(wildcard cmd/examen/*.go)
	go build ./cmd/examen

examensvc.exe: $(wildcard cmd/examensvc/*.go)
	go build ./cmd/examensvc

clean: tidy
	rm setup.exe

tidy:
	rm examen.exe examensvc.exe install.exe 
