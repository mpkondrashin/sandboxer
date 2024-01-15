

.PHONY: clean tidy

setup: install.exe
	go build ./cmd/setup
install: examen.exe examensvc.exe
	go build ./cmd/install
examen:
	go build ./cmd/examen
examensvc:
	go build ./cmd/examensvc

clean: tidy
	rm setup.exe

tidy:
	rm examen.exe examensvc.exe install.exe 
