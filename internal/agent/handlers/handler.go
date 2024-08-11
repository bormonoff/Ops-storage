package handlers

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

var client = resty.New()

func SendPostRequest(url string, headers map[string]string) {
	request := client.R()

	for header, value := range headers {
		request.SetHeader(header, value)
	}

	_, err := request.Post(url)
	if err != nil {
		fmt.Println(err)
	}
}
