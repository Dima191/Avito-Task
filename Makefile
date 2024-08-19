migrate_up:
	migrate -path ./migrations -verbose -database "postgres://{user}:{password}@{host}:{port}/{db_name}?sslmode=disable" up

migrate_down:
	migrate -path ./migrations -verbose -database "postgres://{user}:{password}@{host}:{port}/{db_name}?sslmode=disable" down

build:
	go build ./cmd/app
	./app


.PHONY: migrate_up, migrate_down, build

.DEFAULT_GOAL=build