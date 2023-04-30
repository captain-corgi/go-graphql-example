cd pkg/db/migrations
migrate create -ext sql -dir mysql -seq create_users_table
migrate create -ext sql -dir mysql -seq create_links_table
cd ../../../