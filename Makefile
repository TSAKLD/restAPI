db:
	docker run -v rest:/var/lib/postgresql/data/ --name db -p "5432:5432" -e POSTGRES_PASSWORD=dev -e POSTGRES_USER=kr -e POSTGRES_DB=userdb -d postgres:16.2