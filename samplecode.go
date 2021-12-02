package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	devKey := flag.String("dev-key", "blah", "Dev authentication key")
	pasteCode := flag.String("paste-code", "Hello world, from pastebin!", "Code to be pasted")
	pasteFile := flag.String("paste-file", "", "file containing code")
	pasteOption := flag.String("paste-option", "paste", "Paste option")
	pasteFormat := flag.String("paste-format", "diff", "Format- can be go, c, bash etc")
	pasteExpireDate := flag.String("paste-expiry", "1Y", "Duration of existence")

	flag.Parse()

	pasteUrl := "https://pastebin.com/api/api_post.php"

	var paste string

	if *pasteFile == "" {
		paste = *pasteCode
	} else {
		dat, err := os.ReadFile(*pasteFile)
		check(err)
		paste = string(dat)
	}

	urlValues := map[string][]string{
		"api_dev_key":     {*devKey},
		"api_paste_code":  {paste},
		"api_option":      {*pasteOption},
		"api_format":      {*pasteFormat},
		"api_expiry_date": {*pasteExpireDate},
	}

	resp, err := http.PostForm(pasteUrl, urlValues)

	if err != nil {
		fmt.Errorf(*devKey)
	}

	defer resp.Body.Close()
	link, _ := io.ReadAll(resp.Body)

	fmt.Println(string(link))

}
