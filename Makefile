.PHONY: pg stoppg create-table

pg:
	docker run --name testContainer \
		-e POSTGRES_PASSWORD=secret \
		-e POSTGRES_USER=root \
		-e POSTGRES_DB=test \
		-p 5432:5432 \
		-d postgres:alpine

stoppg:
	docker stop testContainer && docker rm testContainer

create-table:
	migrate -path pkg/database/migrations \
		-database postgres://root:secret@localhost:5432/test?sslmode=disable up
