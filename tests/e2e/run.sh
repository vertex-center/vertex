#!/usr/bin/env bash

echo "Running e2e tests..."

TAG=vertex-e2e-environment

docker build -t $TAG -f tests/e2e/Dockerfile .
docker run -it -d -p 7130:6130 -p 7131:6131 $TAG

while ! nc -z localhost 7130; do sleep 1; done

FILES=$(find . -name "*_e2e_test.go")

if go test -tags e2e -v "$FILES"; then
  EXIT_CODE=0
else
  EXIT_CODE=1
fi

docker stop "$(docker ps -q --filter ancestor=$TAG)"
docker rm "$(docker ps -aq --filter ancestor=$TAG)"

exit "$EXIT_CODE"
