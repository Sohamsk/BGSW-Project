TARGET?=./test.bas

build:
	go build -o ./dist/main main.go

test:
	./dist/main $(TARGET)

clean:
	rm op.txt main
