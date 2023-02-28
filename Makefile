
clean:
	rm -rf resize/resize.zip

zip: clean
	zip -r resize.zip resize/
zip-mv: zip
	mv resize.zip resize