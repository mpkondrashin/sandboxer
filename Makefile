

.PHONY: clean tidy

setup.exe: cmd/setup/embed/install.exe.gz cmd/setup/embed/opengl32.dll.gz $(wildcard cmd/setup/*.go)
	go build ./cmd/setup

cmd/setup/embed/install.exe.gz: install.exe 
	gzip -fc install.exe > cmd/setup/embed/install.exe.gz

cmd/setup/embed/opengl32.dll.gz: resources/opengl32.dll
	gzip -fc resources/opengl32.dll  > cmd/setup/embed/opengl32.dll.gz

install.exe: examen.exe examensvc.exe $(wildcard cmd/install/*.go)
	fyne package --os windows --name ExamenInstaller --appID in.kondrash.examen --appVersion 0.0.1 --icon resource/examen.png --release --sourceDir ./cmd/install
	

examen.exe: $(wildcard cmd/examen/*.go)
	go build ./cmd/examen

examensvc.exe: $(wildcard cmd/examensvc/*.go)
	go build ./cmd/examensvc

clean: tidy
	rm setup.exe

tidy:
	rm examen.exe examensvc.exe install.exe 
