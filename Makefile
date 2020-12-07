version := $(shell if [ ! -z "${TAG}" ]; then echo "${TAG}"; else git describe --tags --always; fi)
date := $(shell date)
export env := ${env}
tag := ${TAG}
branch := $(shell if [ ! -z "${BRANCH}" ]; then echo "${BRANCH}"; else git rev-parse --abbrev-ref HEAD; fi)
builder := $(shell if [ ! -z "${BUILDER}" ]; then echo "${BUILDER}"; elif [ ! -z ${BITBUCKET_BUILD_NUMBER} ]; then git log -1 --pretty=format:'%an' | xargs ; else git config user.name; fi)
ldflags := "-X 'main.version=${version}' -X 'main.branch=${branch}' -X 'main.builder=${builder}' -X 'main.buildDate=${date}'"
org := thenewfork
image-name := api-pipeline

test:
	go test ./... -p 1

build: go-build docker-build

go-build:
	CGO_ENABLED=0 go build -ldflags ${ldflags} -a -installsuffix cgo -o bin/ftp-pipeline ./cmd/ftp-pipeline/...

docker-build:
	docker build --tag registry.unchain.io/${org}/${image-name}:${version} .

push:
	docker push registry.unchain.io/${org}/${image-name}:${version}

version:
	echo ${version}