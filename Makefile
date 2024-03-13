run:
	docker-compose -f docker/docker-compose.yml -f docker/docker-compose.build.yml up --build

run-release:
	docker-compose -f docker/docker-compose.yml up

run-dev:
	docker-compose -f docker/docker-compose.yml -f docker/docker-compose.dev.yml up --build

run-storybook:
	yarn workspace @vertex-center/components storybook

clean:
	docker-compose -f docker/docker-compose.yml down --rmi all
