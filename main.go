package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
)

func main() {
	var (
		outputName string
	)

	flag.StringVar(&outputName, "output", "", "output name")
	flag.StringVar(&outputName, "o", "", "output name")

	flag.Usage = func() {
		fmt.Println("yowes [url] -o my-file")
	}

	flag.Parse()

	args := flag.Args()

	if len(args) <= 0 {
		fmt.Println("url empty or invalid")
		os.Exit(1)
	}

	fmt.Println(args)

	url := args[0]

	fmt.Println(outputName)

	// get response from given url
	response, err := get(url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer response.Body.Close()

	// create output file
	file, err := os.Create(outputName)
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

func isValidURL(url string) bool {
	urlRegex := "^(http:\\/\\/www\\.|https:\\/\\/www\\.|http:\\/\\/|https:\\/\\/)?[a-z0-9]+([\\-\\.]{1}[a-z0-9]+)*\\.[a-z]{2,5}(:[0-9]{1,5})?(\\/.*)?$"
	re := regexp.MustCompile(urlRegex)
	return re.MatchString(url)
}
