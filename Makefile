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
	(cd db/schema && goose ${DB_DRIVER} ${DB_URL} up)

dbmigratedown:
	(cd db/schema && goose ${DB_DRIVER} ${DB_URL} down)

test:
	go test -v -cover ./...

server: 
	go run main.go

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	proto/*.proto

evans: 
	evans --host localhost --port 7777 -r repl

.PHONY: postgrescreate postgresremove dbcreate dbdrop dbconn dbmigrateup dbmigratedown test server proto
