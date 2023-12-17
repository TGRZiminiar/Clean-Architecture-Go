.PHONY: pg stoppg init-table

# Start Pg
pg:
	docker run --name authDB \
		-e POSTGRES_PASSWORD=secret \
		-e POSTGRES_USER=root \
		-e POSTGRES_DB=test \
		-p 5432:5432 \
		-d postgres:alpine

# Stop Postgres db
stoppg:
	docker stop authDB && docker rm authDB

# Init table
init-db:
	migrate -path pkg/database/migrations \
		-database postgres://mix:secret@localhost:5432/test?sslmode=disable up

# Init migrate
# Sample: make init-mg name=create
# Complete Command: migrate create -ext sql -dir pkg/database/migrations -seq init_schema2
create-mg:
	migrate create -ext sql -dir pkg/database/migrations -seq $(name)

u:
	go run main.go ./env/dev/.env.dev
