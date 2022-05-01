package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// shortened for readability
type (
	Req events.LambdaFunctionURLRequest
	Res events.LambdaFunctionURLResponse
)

func handle(_ context.Context, request *Req) (response *Res, err error) {

	if !strings.HasPrefix(request.RawPath, "/mlctrez") {
		response = &Res{StatusCode: http.StatusNotFound}
		return
	}

	dumpRequest(request)

	urlString := fmt.Sprintf("https://goreportcard.com/badge/github.com%s", request.RawPath)

	var u *url.URL
	if u, err = u.Parse(urlString); err != nil {
		response = &Res{StatusCode: http.StatusInternalServerError}
		return response, err
	}

	req := &http.Request{URL: u}

	var res *http.Response

	if res, err = http.DefaultClient.Do(req); err != nil {
		response = &Res{StatusCode: http.StatusInternalServerError}
		return response, err
	}

	defer func() { _ = res.Body.Close() }()

	var resBody []byte

	if resBody, err = io.ReadAll(res.Body); err != nil {
		response = &Res{StatusCode: http.StatusInternalServerError}
		return response, err
	}

	response = &Res{
		StatusCode: res.StatusCode,
		Headers: map[string]string{
			"Content-Type":  res.Header.Get("Content-Type"),
			"Cache-Control": "max-age=3600",
		},
		Body: string(resBody),
	}

	return response, nil
}

func dumpRequest(request *Req) {
	if marshal, err := json.Marshal(request); err == nil {
		fmt.Println(string(marshal))
	}
}

func main() {
	lambda.Start(handle)
}
