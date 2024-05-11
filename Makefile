db:
	docker rm -f db
	docker run -v rest:/var/lib/postgresql/data/ --name db -p "5432:5432" --restart=always -e POSTGRES_PASSWORD=dev -e POSTGRES_USER=kr -e POSTGRES_DB=userdb -d postgres:16.2

migrate-up:
	 goose -dir migrations postgres "user=kr host=localhost port=5432 password=dev dbname=userdb sslmode=disable"  up

migrate-down:
	 goose -dir migrations postgres "user=kr host=localhost port=5432 password=dev dbname=userdb sslmode=disable"  down

 migrate-reset:
	 goose -dir migrations postgres "user=kr host=localhost port=5432 password=dev dbname=userdb sslmode=disable"  reset && \
 	 goose -dir migrations postgres "user=kr host=localhost port=5432 password=dev dbname=userdb sslmode=disable"  up

run-tests:
	docker rm -f db_test || true
	docker run --name --rm db_test -e TZ=UTC -p "5433:5432" -e POSTGRES_PASSWORD=postgres -d postgres:16.2 && \
	timeout 3 && \
	goose -dir migrations postgres "user=postgres host=localhost port=5433 password=postgres dbname=postgres sslmode=disable"  up && \
	go test ./...

kafka:
	docker rm -f kafka
	docker run --name kafka --restart=always -d -p 9092:9092 apache/kafka:3.7.0

redis:
	docker rm -f redis
	docker run --name redis --restart=always -d -p 6379:6379 redis:7.2