e2e:
	./tests/e2e/run.sh

build-docker-microservices:
	COMPOSE_PARALLEL_LIMIT=1 docker-compose -f docker/micro.docker-compose.yml --project-name vertex build

build-docker-bundle:
	COMPOSE_PARALLEL_LIMIT=1 docker-compose -f docker/bundle.docker-compose.yml --project-name vertex build

start-docker-microservices:
	COMPOSE_PARALLEL_LIMIT=1 docker-compose -f docker/micro.docker-compose.yml --project-name vertex up

start-docker-bundle:
	COMPOSE_PARALLEL_LIMIT=1 docker-compose -f docker/bundle.docker-compose.yml --project-name vertex up

stop-docker-microservices:
	COMPOSE_PARALLEL_LIMIT=1 docker-compose -f docker/micro.docker-compose.yml --project-name vertex down

stop-docker-bundle:
	COMPOSE_PARALLEL_LIMIT=1 docker-compose -f docker/bundle.docker-compose.yml --project-name vertex down
