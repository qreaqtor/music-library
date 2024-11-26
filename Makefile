PG_URL=postgres://user:password@localhost:5432/songs?sslmode=disable

.PHONY: up-dev
up-dev:
	export CONFIG_PATH="./config/dev.env" && docker-compose up -d

.PHONY: up-local
up-local:
	docker-compose up -d postgres migrations

.PHONY: down
down:
	docker-compose down && docker rmi music-library-app music-library-migrations

.PHONY: .run-app-local
run-app-local:
	export CONFIG_PATH="./config/local.env" && go run cmd/music-library/main.go

.PHONY: .gen-swagger
gen-swagger:
	swag init -g internal/api/songs.go

.PHONY: .migration-up
migration-up:
	$(eval PG_URL?=$(PG_URL))
	goose -dir ./migrations postgres "$(PG_URL)" up

.PHONY: .migration-down
migration-down:
	$(eval PG_URL?=$(PG_URL))
	goose -dir ./migrations postgres "$(PG_URL)" down

.PHONY: .migration-status
migration-status:
	$(eval PG_URL?=$(PG_URL))
	goose -dir ./migrations postgres "$(PG_URL)" status

.PHONY: .migration-create-sql
migration-create-sql:
	goose -dir ./migrations create $(filter-out $@,$(MAKECMDGOALS)) sql
