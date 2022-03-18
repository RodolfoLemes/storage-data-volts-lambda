build:
	GOOS=linux go build -o bin/main
	rm -r -f storage-data-volts-lambda
	zip storage-data-volts-lambda bin/main