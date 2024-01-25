BINARY_NAME=dvault

build:
	go build -o ${BINARY_NAME} main.go

run-dev:
	go run main.go

run: build
	./dvault

clean:
	go clean
	rm ${BINARY_NAME}

vendor:
	go mod vendor

install-migration-tool:
	brew install golang-migrate

# example: migration-create name=init_schema
migration-create:
	migrate create -ext sql -dir db/migrations -seq ${name}

migration-up:
	migrate -path db/migrations/ -database "postgresql://postgres:postgres@localhost:5432/dvault?sslmode=disable" -verbose up

migration-down:
	migrate -path db/migrations/ -database "postgresql://postgres:postgres@localhost:5432/dvault?sslmode=disable" -verbose down
