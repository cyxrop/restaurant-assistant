ACCESS_SECRET = some_secret_string

# Database configuration
RA_DB_DRIVER ?= postgres
RA_DB_DRIVER_VERSION ?= 10.5
RA_DB_NAME ?= restaurant_assistant
RA_DB_USERNAME ?= postgres
RA_DB_PASSWORD ?= postgres
RA_DB_LOCATION ?= ~/docker/postgres
RA_DB_PATH ?= /var/lib/postgresql/data
RA_DB_HOST ?= postgres
RA_DB_PORT ?= 5432

run:
	@echo "+ $@"
	ACCESS_SECRET=$(ACCESS_SECRET) \
	go run cmd/cmd.go;


HAS_DB_RUNNED := $(shell docker ps | grep $(RA_DB_HOST))
HAS_DB_EXITED := $(shell docker ps -a | grep $(RA_DB_HOST))

db-up:
	@echo "+ $@"
ifndef HAS_DB_RUNNED
ifndef HAS_DB_EXITED
	@mkdir -p $(RA_DB_LOCATION)
	@docker run -d	--name $(RA_DB_DRIVER) \
	-p $(RA_DB_PORT):$(RA_DB_PORT) \
	-e "POSTGRES_DB=$(RA_DB_NAME)" \
	-e "POSTGRES_USER=$(RA_DB_USERNAME)" \
	-e "POSTGRES_PASSWORD=$(RA_DB_PASSWORD)" \
	-v $(RA_DB_LOCATION):$(RA_DB_PATH) \
	$(RA_DB_DRIVER):$(RA_DB_DRIVER_VERSION)
	@sleep 45
else
	@docker start $(RA_DB_HOST)
endif
endif

db-down:
	@echo "+ $@"
ifdef HAS_DB_RUNNED
	@docker stop $(RA_DB_HOST)
endif

pg:
	@echo "+ $@"
	@docker exec -it postgres bash -c "psql -U postgres restaurant_assistant"

migrate-up:
	@echo "+ $@"
	@sql-migrate up -env="development" \
	-config=db/config/dbconfig.yml

migrate-down:
	@echo "+ $@"
	@sql-migrate down -env="development" \
	-config=db/config/dbconfig.yml

clear-db: migrate-down migrate-up
	@echo "+ $@"

db-models:
	@echo "+ $@"
	@sqlboiler -c=./.sqlboiler.toml -d=true psql

up: db-up
	@echo "+ $@"
	@docker-compose up -d

down: db-down
	@echo "+ $@"
	@docker-compose down

.PHONY: run \
		pg \
		migrate-up \
		migrate-down \
		db-models \
		clear-db \
		up \
		down \
		db-up