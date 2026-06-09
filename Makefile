all:
	cd gob && go run db_builder.go
	cd gob && go run server.go