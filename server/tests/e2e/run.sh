#!/usr/bin/env bash

docker-compose -f docker/bundle.docker-compose.yml -f docker/bundle.docker-compose.build.yml up -d --build

#while ! nc -z localhost 7130; do sleep 1; done

echo "Running e2e tests..."

if venom run "**/*_test.yml"; then
    EXIT_CODE=0
else
    EXIT_CODE=1
fi

docker-compose -f docker/bundle.docker-compose.yml -f docker/bundle.docker-compose.build.yml down

exit "$EXIT_CODE"
