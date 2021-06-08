build:
	go build -ldflags "-s -w" -o bin/

modvendor:
	go build -mod=vendor -o bin/

clean:
	rm -rf ./bin
