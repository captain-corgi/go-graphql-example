app:
	go run cmd/app/main.go

hackernews:
	go run cmd/hackernews/main.go

tidy:
	go mod tidy
	go mod vendor

gqlgen:
	@cd internal/hackernews
	gqlgen.exe $1
	@cd ../../

mysqlconnect:
	mysql -u root -p

mysqlcreate:
	@cd pkg/db/migrations
	migrate create -ext sql -dir mysql -seq create_users_table
	migrate create -ext sql -dir mysql -seq create_links_table
	@cd ../../../

mysqlrun:
	docker run -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=dbpass -e MYSQL_DATABASE=hackernews -d mysql:latest