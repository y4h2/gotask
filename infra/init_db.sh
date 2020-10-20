#!/bin/sh

ROOT_DIR=$(git rev-parse --show-toplevel)
CONTAINER_NAME=some-postgres
docker rm -f ${CONTAINER_NAME}
docker run -d \
    --name ${CONTAINER_NAME} \
    -p 5432:5432 \
    -e POSTGRES_USER=myuser \
    -e POSTGRES_PASSWORD=mypassword \
    -v ${ROOT_DIR}/infra/data/:/docker-entrypoint-initdb.d/ \
    postgres