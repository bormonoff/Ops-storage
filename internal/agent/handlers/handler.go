package handlers

import (
	"github.com/go-resty/resty/v2"
)

var client = resty.New()

func SendPostRequest(url string, headers map[string]string, body any) error {
	request := client.R().SetBody(body)

	for header, value := range headers {
		request.SetHeader(header, value)
	}

	_, err := request.Post(url)
	if err != nil {
		return err
	}

	return nil
}
