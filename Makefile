.PHONY: tidy build-db build-server

tidy:
	cd gob && go mod tidy

build-db: tidy
	cd gob && go run db_builder.go

build-server: tidy
	cd gob && go build -o server .