TARGET?=./test.bas
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
	cat op.json

.PHONY: clean
clean:
	rm -rf op.json ./dist
