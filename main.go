package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
)

const (
	// Version of yowes
	Version = "1.0.0"

	// Banner of yowes
	Banner = `                                          
		
	▓██   ██▓ ▒█████   █     █░▓█████   ██████ 
	▒██  ██▒▒██▒  ██▒▓█░ █ ░█░▓█   ▀ ▒██    ▒ 
	▒██ ██░▒██░  ██▒▒█░ █ ░█ ▒███   ░ ▓██▄   
	░ ▐██▓░▒██   ██░░█░ █ ░█ ▒▓█  ▄   ▒   ██▒
	░ ██▒▓░░ ████▓▒░░░██▒██▓ ░▒████▒▒██████▒▒
	██▒▒▒ ░ ▒░▒░▒░ ░ ▓░▒ ▒  ░░ ▒░ ░▒ ▒▓▒ ▒ ░
	▓██ ░▒░   ░ ▒ ▒░   ▒ ░ ░   ░ ░  ░░ ░▒  ░ ░
	▒ ▒ ░░  ░ ░ ░ ▒    ░   ░     ░   ░  ░  ░  
	░ ░         ░ ░      ░       ░  ░      ░  
	░ ░                                       
												
   `
)

func main() {
	var (
		showVersion bool
	)

	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.BoolVar(&showVersion, "v", false, "show version")

	flag.Usage = func() {
		fmt.Println()
		fmt.Println(Banner)
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
		fmt.Println(Banner)
		fmt.Printf("  yowes version %s\n", Version)
		fmt.Println()
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

	// measure execution
	done := make(chan bool, 1)
	kill := make(chan os.Signal, 1)

	// tick every 500 millisecond
	ticker := time.NewTicker(time.Millisecond * 500)
	fmt.Print("please wait, ")

	// notify when user interrupt the process
	signal.Notify(kill, syscall.SIGINT, syscall.SIGTERM)

	go waitOSNotify(kill, done)

	// show execution process
	go measureExecution(done, ticker)

	// get response from given url
	response, err := httpGet(url)
	if err != nil {
		done <- true
		fmt.Println("cannot perform a request")
		os.Exit(1)
	}

	defer response.Body.Close()

	// check status code
	if response.StatusCode != 200 {
		done <- true
		fmt.Printf("request fail, status = %d", response.StatusCode)
		os.Exit(1)
	}

	fileName := getFileName(url)

	// create output file
	file, err := os.Create(fileName)
	if err != nil {
		done <- true
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	err = readWrite(response.Body, file, done)
	if err != nil {
		done <- true
		fmt.Println(err)
		os.Exit(1)
	}

}

func waitOSNotify(kill chan os.Signal, done chan bool) {
	for {
		select {
		case <-kill:
			fmt.Println("download interrupted")
			done <- true
			return
		}
	}
}

func measureExecution(done chan bool, ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			fmt.Print(".")
		case <-done:
			fmt.Println()
			os.Exit(0)
			return
		}
	}
}

func httpGet(url string) (*http.Response, error) {
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
		IdleConnTimeout:     10 * time.Second,
	}

	httpClient := &http.Client{
		//Timeout:   time.Second * 10,
		Transport: transport,
	}
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

func readWrite(in io.Reader, out io.Writer, done chan bool) error {
	buffer := make([]byte, 1024)
	reader := bufio.NewReader(in)

	for {
		line, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		// write
		_, err = out.Write(buffer[:line])
		if err != nil {
			return err
		}
	}

	done <- true

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
