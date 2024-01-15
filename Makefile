

.PHONY: clean tidy

setup.exe: install.exe.gz opengl32.dll.gz
	go build ./cmd/setup

install.exe.gz: install.exe
	gzip -f install.exe

opengl32.dll.gz: resources/opengl32.dll
	gzip -fc resources/opengl32.dll  > opengl32.dll.gz

install.exe: examen.exe examensvc.exe
	go build ./cmd/install

examen.exe:
	go build ./cmd/examen

examensvc.exe:
	go build ./cmd/examensvc

clean: tidy
	rm setup.exe

tidy:
	rm examen.exe examensvc.exe install.exe 
