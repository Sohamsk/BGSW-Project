TARGET?=./testfiles/property_test.cls
.PHONY: setup
setup:
	curl https://www.antlr.org/download/antlr-4.13.2-complete.jar -o ./parser/antlr.jar
	go generate .\...

.PHONY: build
build:
	go build -o ./dist/main.exe .

.PHONY: test
test: build
	./dist/main.exe $(TARGET)
	

.PHONY: clean
clean:
	rm -r .\output
	rm -r .\dist
