TARGET?=./testfiles/funcTest.cls
.PHONY: setup
setup:
	curl https://www.antlr.org/download/antlr-4.13.2-complete.jar -o ./parser/antlr.jar
	go generate ./...

.PHONY: build
build:
	go build -o ./dist/main main.go

.PHONY: test
test: build
	./dist/main $(TARGET)

.PHONY: clean
clean:
	rm -rf ./output ./dist
