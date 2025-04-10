build: 
	@go build -o bin/dmfs

run: build
	@./bin/dmfs

test:
	@go test -v ./...

clean:
	@rm -f bin/dmfs

