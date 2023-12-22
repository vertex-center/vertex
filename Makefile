e2e:
	./tests/e2e/run.sh

docker-compose-build:
	COMPOSE_PARALLEL_LIMIT=1 docker-compose -f docker/micro.docker-compose.yml --project-name vertex build

docker-compose-up:
	COMPOSE_PARALLEL_LIMIT=1 docker-compose -f docker/micro.docker-compose.yml --project-name vertex up
