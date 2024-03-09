package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	bd "github.com/lutzky/go-bidi"
)

// aa123aa1234aa1aaa
// aa321aa4321aa1aaa

func isHebrew(r rune) bool {
	if unicode.Is(unicode.Hebrew, r) {
		return true
	}
	return false
	if r == ' ' {
		return true
	}
	return false
}

func ReverseHebrew(s string) string {
	runes := []rune(s)
	state := 0
	i, j := 0, 0
	for {
		if i == len(runes) {
			return string(runes)
		}
		switch state {
		case 0: // seek for Hebrew
			if isHebrew(runes[i]) {
				j = i + 1
				state = 1
				continue
			}
			i++
			if i == len(runes) {
				return string(runes)
			}
		case 1: // seek for non hebrew
			if j == len(runes) || !isHebrew(runes[j]) {
				state = 0
				for q, p := i, j-1; q < p; q, p = q+1, p-1 {
					runes[q], runes[p] = runes[p], runes[q]
				}
				i = j
				continue
			} else {
				j++
			}
		}
	}

}

func reverse(s string) string {
	parts := strings.Split(s, "_")
	reversed := make([]string, len(parts))
	for i := 0; i < len(parts); i++ {
		reversed[i] = ReverseHebrew(parts[len(parts)-i-1])
	}
	return strings.Join(reversed, "_")
}

func process(s string) string {
	parts := strings.Split(s, ".")
	for i := 0; i < len(parts); i++ {
		parts[i] = reverse(parts[i])
	}
	return strings.Join(parts, ".")
}

func split(s string) (result []string) {
	if s == "" {
		return nil
	}
	runes := []rune(s)
	h1 := unicode.Is(unicode.Hebrew, runes[0])
	word := string(runes[0])
	for i := 1; i < len(runes); i++ {
		r := runes[i]
		h2 := unicode.Is(unicode.Hebrew, r)
		if h1 != h2 {
			result = append(result, word)
			word = ""
			h1 = h2
		}
		word += string(r)
	}
	if word != "" {
		result = append(result, word)
	}
	return
}

func main() {
	fontFileName := "DroidSansHebrew-Regular.ttf"
	os.Setenv("FYNE_FONT", filepath.Join("../../../resources", fontFileName))

	s := "ארכיון_חשבוניות_10_2020 עד_10_2023_7921117.csv"
	//s1 := bidi.ReverseString(s)
	s1, err := bd.Display(s)
	log.Println(err)
	a := app.New()
	w := a.NewWindow("Hello World")
	s1L := widget.NewLabel(s1)
	s2L := widget.NewLabel(s)
	b := widget.NewButton("Notification", func() {
		//xplatform.Alert(globals.AppID, "Sandboxer", "High risk treat detected", "abcd.exe")
		//notif := fyne.NewNotification("Title", "content")
		//err := beeep.Notify("Title", "Message 1 body", "../../../resources/icon_transparent.png")
		//if err != nil {
		//panic(err)
		//}
		//err := Alert("Title", "Message body", "../../../resources/icon_transparent.png")
		//if err != nil {
		//	panic(err)
		//}
		//a.SendNotification(notif)
	})
	vbox := container.NewVBox(s1L, s2L, b)
	w.SetContent(vbox)
	w.ShowAndRun()
	return
	/*
		_ = "ארכיון_חשבוניות_10_2020 עד_10_2023_7921117.csv"
		s := "ארכיון7921117.csv"
		fmt.Println([]rune(s))
		fmt.Println()
		fmt.Println("Just print:")
		fmt.Println(s)
		fmt.Println()
		fmt.Println("hex:")
		for i := 0; i < len(s); i++ {
			fmt.Printf("%x ", s[i])
		}
		fmt.Println()
		d := split(s)
		fmt.Println(strings.Join(d, "|"))
		fmt.Println(d)
		return
		//s := "שכירות.pdf"
		fmt.Println()
		fmt.Println("Just print:")
		fmt.Println(s)
		sr := ReverseHebrew(s)
		fmt.Println()
		fmt.Println("ReverseHebrew:")
		fmt.Println(sr)

		srr := process(s)
		fmt.Println()
		fmt.Println("ReverseHebrew:")
		fmt.Println(srr)

		fmt.Println()
		fmt.Println("print char by char:")
		for _, r := range s {
			fmt.Printf("%s", string(r))
		}
		fmt.Println()

		fmt.Println()
		fmt.Println("hex:")
		for i := 0; i < len(s); i++ {
			fmt.Printf("%x ", s[i])

		}
		fmt.Println()

		fontFileName := "DroidSansHebrew-Regular.ttf"
		os.Setenv("FYNE_FONT", filepath.Join("../../../resources", fontFileName))

		a := app.New()
		w := a.NewWindow("Hello World")
		s1 := widget.NewLabel(s)
		s2 := widget.NewLabel(sr)
		s3 := widget.NewLabel(srr)
		vbox := container.NewVBox(s1, s2, s3)
		w.SetContent(vbox)
		w.ShowAndRun()*/
}
