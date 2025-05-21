include ./Makefile.variable
# Here we import the variables for our makefiles. Please check them out if you are unsure what they do/are

# Create postgres container
postgrescreate:
	docker run --name ${POSTGRES_CONTAINER_NAME} -p ${DB_INTERNAL_PORT}:${DB_EXTERNAL_PORT} -e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_USER_PASSWORD} -d postgres

# Create the Database
dbcreate:
	docker exec -it ${POSTGRES_CONTAINER_NAME} createdb --username=${DB_USER} --owner=${DB_USER} ${DB_NAME}

# Test connection
dbconn:
	psql ${DB_URL}

# Drop the database
dbdrop:
	docker exec -it ${POSTGRES_CONTAINER_NAME} dropdb ${DB_NAME}

postgresremove:
	docker stop ${POSTGRES_CONTAINER_NAME}
	docker rm ${POSTGRES_CONTAINER_NAME}

dbmigrateup:
	(cd internal/db/schema && goose ${DB_DRIVER} ${DB_URL} up)

dbmigratedown:
	(cd internal/db/schema && goose ${DB_DRIVER} ${DB_URL} down)

test:
	go test -v -cover ./...

server: 
	go run cmd/gonuxt-starter/main.go

proto:
	rm -f internal/pb/*.go
	rm -f api/swagger/*.swagger.json
	protoc --proto_path=internal/proto --go_out=internal/pb --go_opt=paths=source_relative \
	--go-grpc_out=internal/pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=internal/pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=api/swagger --openapiv2_opt=allow_merge=true,merge_file_name=gonuxt,json_names_for_fields=false \
	--experimental_allow_proto3_optional \
	internal/proto/*.proto
	(cd web && npm run api:generate)

evans: 
	evans --host localhost --port 7777 -r repl

tools:
	@echo "Installing tools from $(TOOLS_FILE)..."
	@grep '_ "' $(TOOLS_FILE) | while read line; do \
		tool=$$(echo $$line | cut -d'"' -f2); \
		go install $$tool; \
	done

build: 
	go build -o ./bin/gonuxt-starter ./cmd/gonuxt-starter

godeps:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest



.PHONY: postgrescreate postgresremove dbcreate dbdrop dbconn dbmigrateup dbmigratedown test server proto tools
