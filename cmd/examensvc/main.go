package main

func InspectFile(filePath string) {

}

func main() {
	var inbox StringChannel
	go SubmitDispatch(inbox)
	for {
		s := <-inbox
		go InspectFile(s)
	}
}
