
build:
	@mkdir -p temp
	@GOOS=linux CGO_ENABLED=0 go build -o temp/main main.go

zip: build
	@cd temp && zip -q function.zip main

upload: zip
	@aws lambda update-function-code \
		--function-name arn:aws:lambda:us-east-1:359625541351:function:badge \
		--zip-file fileb://temp/function.zip --publish > /dev/null

clean:
	@rm -rf temp