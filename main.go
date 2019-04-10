package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	// Version of yowes
	Version = "0.0.0"
)

func main() {
	var (
		showVersion bool
	)

	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.BoolVar(&showVersion, "v", false, "show version")

	flag.Usage = func() {
		fmt.Println()
		fmt.Println("Usage: ")
		fmt.Println("yowes [url]")
		fmt.Println()
		fmt.Println("-h | -help (show help)")
		fmt.Println("-v | -version (show version)")
		fmt.Println("---------------------------")
		fmt.Println()
	}

	flag.Parse()

	if showVersion {
		fmt.Printf("yowes version %s\n", Version)
		os.Exit(0)
	}

	args := flag.Args()

	if len(args) <= 0 {
		fmt.Println("url empty or invalid")
		os.Exit(1)
	}

	url := args[0]

	if !isValidURL(url) {
		fmt.Println("url empty or invalid")
		os.Exit(1)
	}

	// get response from given url
	response, err := get(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer response.Body.Close()

	fileName := getFileName(url)

	// create output file
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	read(response.Body, file)

}

func get(url string) (*http.Response, error) {
	httpClient := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func read(in io.Reader, out io.Writer) error {
	body, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	_, err = out.Write(body)
	if err != nil {
		return err
	}

	return nil
}

func getFileName(urlParam string) string {
	urls := strings.Split(urlParam, "/")
	fileName := urls[len(urls)-1]

	return fileName
}

func isValidURL(urlParam string) bool {
	urlRegex := "^(http:\\/\\/www\\.|https:\\/\\/www\\.|http:\\/\\/|https:\\/\\/)?[a-z0-9]+([\\-\\.]{1}[a-z0-9]+)*\\.[a-z]{2,5}(:[0-9]{1,5})?(\\/.*)?$"
	re := regexp.MustCompile(urlRegex)
	return re.MatchString(urlParam)
}
