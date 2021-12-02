package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const pasteUrl = "https://pastebin.com/api/api_post.php"

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// return errors, caller should handle what to do with the returned error

func getPasteCode(pasteCode *string, pasteFile *string, pasteClip *bool) string {
	if *pasteClip {
		return getClipContent()
	}

	if *pasteFile == "" {
		return *pasteCode
	} else {
		dat, err := os.ReadFile(*pasteFile)
		check(err)
		return string(dat)
	}
}

// return errors, caller should handle what to do with the returned error
func getClipContent() string {
	cmd := exec.Command("xclip", "-o")

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()
	check(err)

	if errb.String() != "" {
		log.Fatalf("error using xclip: %+v", errb.String())
	}

	return outb.String()
}

func main() {
	devKey := flag.String("dev-key", "", "Dev authentication key")
	pasteCode := flag.String("paste-code", "Hello world, from pastebin!", "Code to be pasted")
	pasteFile := flag.String("paste-file", "", "file containing code")
	pasteOption := flag.String("paste-option", "paste", "Paste option")
	pasteFormat := flag.String("paste-format", "diff", "Format- can be Go, C, Bash etc")
	pasteExpireDate := flag.String("paste-expiry", "10M", "Duration of existence")
	pasteClip := flag.Bool("c", false, "Paste content from clipboard")

	flag.Parse()

	if *devKey == "" {
		log.Fatalf("\nYou forgot to give me your developer key.\nSet it using the '-dev-key' flag")
	}

	// move this code to a separate function called for e.g. makeRequest
	paste := getPasteCode(pasteCode, pasteFile, pasteClip)

	urlValues := map[string][]string{
		"api_dev_key":           {*devKey},
		"api_paste_code":        {paste},
		"api_option":            {*pasteOption},
		"api_paste_format":      {*pasteFormat},
		"api_paste_expire_date": {*pasteExpireDate},
	}

	resp, err := http.PostForm(pasteUrl, urlValues)
	check(err)

	defer resp.Body.Close()

	// error handling missing
	pasteLink, _ := io.ReadAll(resp.Body)
	fmt.Println(string(pasteLink))
}
