TARGET?=./testfiles/withTest.cls
.PHONY: setup
setup:
	curl https://www.antlr.org/download/antlr-4.13.2-complete.jar -o ./parser/antlr.jar
	go generate ./...

.PHONY: build
build:
	go build -o ./dist/main .

.PHONY: test
test: build
	./dist/main $(TARGET)
	@echo "**logs**"
	@cat ./output/logs.log

.PHONY: clean
clean:
	rm -rf ./output ./dist
