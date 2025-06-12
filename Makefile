up:
	docker compose --env-file /dev/null up --build

up-standalone:
	docker compose -f compose-standalone.yml --env-file /dev/null up --build
