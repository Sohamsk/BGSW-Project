TARGET?=./test.bas

setup:
	curl https://www.antlr.org/download/antlr-4.13.2-complete.jar -o ./parser/antlr.jar
	go generate ./...

build:
	go build -o ./dist/main main.go

test: build
	./dist/main $(TARGET)

clean:
	rm -rf op.json ./dist
