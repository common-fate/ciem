PREFIX?=/usr/local

go-binary:
	go build -o ./bin/cf cmd/cli/main.go

cli: go-binary
	mv ./bin/cf ${PREFIX}/bin/
	chmod +x ${PREFIX}/bin/cf

clean:
	rm ${PREFIX}/bin/cf