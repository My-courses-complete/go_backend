postgres:
	docker run --name postgres -e POSTGRES_PASSWORD=password -d -p 5432:5432 postgres:12-alpine

create-db:
	docker exec -it postgres createdb --username=postgres --owner=postgres go_course_bank

drop-db:
	docker exec -it postgres dropdb --username=postgres go_course_bank

migration-up:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/go_course_bank?sslmode=disable" -verbose up

migration-down:
	migrate -path db/migration -database "postgresql://postgres:password@localhost:5432/go_course_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: postgres create-db drop-db migration-up migration-down
