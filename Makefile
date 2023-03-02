
clean:
	rm -rf function.zip

npminstall:
	npm --prefix ./resize install

build: clean npminstall
	
