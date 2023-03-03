
clean:
	rm -rf build
	rm -rf main

build.linux: clean
	GOOS=linux GOARCH=amd64 go build -o main main.go 
	mkdir -p build 
	mv main build	

build: clean
	go build -o main main.go 
	mkdir -p build 
	mv main build	
