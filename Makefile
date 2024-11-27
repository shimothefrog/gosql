# dump to local init.sql file.
dump-db:
	docker compose exec db pg_dump -U admin --inserts sample > ./db/init.sql