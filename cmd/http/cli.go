package main

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"

	bunnyshieldm2m "github.com/le0developer/bunnyshield-m2m"
)

func main() {
	if len(os.Args) != 2 {
		println("Usage: cli <url>")
		os.Exit(1)
	}

	jar, _ := cookiejar.New(nil)
	client := http.Client{
		Transport: &customTransport{
			headers: http.Header{
				"user-agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.3"},
			},
			base: &http.Transport{},
		},
		Jar: jar,
	}

	url := os.Args[1]
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		println("Failed to create request: ", err)
		os.Exit(1)
	}

	res, err := client.Do(req)
	if err != nil {
		println("Failed to perform request: ", err)
		os.Exit(1)
	}

	cookies, err := bunnyshieldm2m.SolveResponse(res, bunnyshieldm2m.HTTPSolverConfig{
		SolveConfig: bunnyshieldm2m.DefaultSolveConfig(),
		Client:      client,
	})
	if err != nil {
		println("Failed to solve response: ", err)
		os.Exit(1)
	}

	for _, cookie := range cookies {
		println(cookie.String())
	}

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		println("Failed to create request: ", err)
		os.Exit(1)
	}

	finalRes, err := client.Do(req)
	if err != nil {
		println("Failed to perform request: ", err)
		os.Exit(1)
	}

	content, err := io.ReadAll(finalRes.Body)
	if err != nil {
		println("Failed to read response: ", err)
		os.Exit(1)
	}

	println(string(content))
}
