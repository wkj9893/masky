package main

import (
	"fmt"
	"net/http"
)

func main() {
	client := http.Client{Transport: &http.Transport{
		Proxy: nil,
	}}
	resp, _ := client.Get("http://example.com")

	fmt.Println(resp)
}
