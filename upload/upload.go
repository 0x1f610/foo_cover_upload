package upload

import (
	"bytes"
	"fmt"
	"io"

	"net/http"
	"os"
)

func Run(url string, auth string) {
	// STDIN receives a file path!
	stdin, err := io.ReadAll(os.Stdin)

	if err != nil {
		panic(err)
	}

	fileData, err := os.ReadFile(string(stdin))
	if err != nil {
		os.Exit(1)
	}

	loadedImage := bytes.NewReader(fileData)

	req, err := http.NewRequest("POST", url, loadedImage)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("User-Agent", "foo_cover_upload/0.1")
	if auth != "" {
		req.Header.Add("Authorization", auth)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		fmt.Println()
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
	if resp.StatusCode != 200 {
		os.Exit(1)
	}
}
