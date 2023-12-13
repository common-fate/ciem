PREFIX?=/usr/local

go-binary:
	go build -o ./bin/dcf cmd/cli/main.go

cli: go-binary
	mv ./bin/dcf ${PREFIX}/bin/
	chmod +x ${PREFIX}/bin/dcf

clean:
	rm ${PREFIX}/bin/dcf