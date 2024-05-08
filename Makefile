db:
	docker run -v rest:/var/lib/postgresql/data/ --name db -p "5432:5432" --restart always -e POSTGRES_PASSWORD=dev -e POSTGRES_USER=kr -e POSTGRES_DB=userdb -d postgres:16.2

migrate-up:
	 goose -dir migrations postgres "user=kr host=localhost port=5432 password=dev dbname=userdb sslmode=disable"  up

migrate-down:
	 goose -dir migrations postgres "user=kr host=localhost port=5432 password=dev dbname=userdb sslmode=disable"  down

run-tests:
	docker rm -f db_test || true
	docker run --name db_test -e TZ=UTC -p "5433:5432" -e POSTGRES_PASSWORD=postgres -d postgres:16.2 && \
	timeout 3 && \
	goose -dir migrations postgres "user=postgres host=localhost port=5433 password=postgres dbname=postgres sslmode=disable"  up && \
	go test ./...