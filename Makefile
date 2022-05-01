
build:
	@mkdir -p temp
	@GOOS=linux CGO_ENABLED=0 go build -o temp/main main.go

zip: build
	@cd temp && zip -q function.zip main

upload: zip
	@aws lambda update-function-code \
		--function-name arn:aws:lambda:us-east-1:359625541351:function:badge \
		--zip-file fileb://temp/function.zip --publish > /dev/null

lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.45.2 golangci-lint run -v

clean:
	@rm -rf temp