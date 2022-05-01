package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func handle(ctx context.Context, request events.LambdaFunctionURLRequest) (response events.LambdaFunctionURLResponse, err error) {

	if !strings.HasPrefix(request.RawPath, "/mlctrez") {
		response = events.LambdaFunctionURLResponse{StatusCode: http.StatusNotFound}
		return
	}

	dumpRequest(request)

	urlString := fmt.Sprintf("https://goreportcard.com/badge/github.com%s", request.RawPath)

	var u *url.URL
	if u, err = u.Parse(urlString); err != nil {
		response = events.LambdaFunctionURLResponse{StatusCode: http.StatusInternalServerError}
		return response, err
	}

	req := &http.Request{URL: u}

	var res *http.Response

	if res, err = http.DefaultClient.Do(req); err != nil {
		response = events.LambdaFunctionURLResponse{StatusCode: http.StatusInternalServerError}
		return response, err
	}

	var resBody []byte

	if resBody, err = ioutil.ReadAll(res.Body); err != nil {
		response = events.LambdaFunctionURLResponse{StatusCode: http.StatusInternalServerError}
		return response, err
	}

	response = events.LambdaFunctionURLResponse{
		StatusCode: res.StatusCode,
		Headers: map[string]string{
			"Content-Type":  res.Header.Get("Content-Type"),
			"Cache-Control": "max-age=3600",
		},
		Body: string(resBody),
	}

	return
}

func dumpRequest(request events.LambdaFunctionURLRequest) {
	if marshal, err := json.Marshal(request); err == nil {
		fmt.Println(string(marshal))
	}
}

func main() {
	lambda.Start(handle)
}
