PG_URL=postgres://user:password@localhost:5432/songs?sslmode=disable

.PHONY: up-dev
up:
	docker-compose up -d

.PHONY: up-local
up-local:
	docker-compose up -d postgres migrations

.PHONY: down
down:
	docker-compose down && docker rmi music-library-app

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
