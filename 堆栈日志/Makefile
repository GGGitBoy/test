NAME = ryan/airflow
VERSION = 1.10.4

.PHONY: build start push

build:  build-version

build-version:
        docker build -t ${NAME}:${VERSION}  .

tag-latest:
        docker tag ${NAME}:${VERSION} ${NAME}:latest

start:
        docker run -it --rm ${NAME}:${VERSION} /bin/bash

push:   build-version tag-latest
        docker push ${NAME}:${VERSION}; docker push ${NAME}:latest