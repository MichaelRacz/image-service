# image-service

## Architecture

The service consists of 3 components
- **queue**: stores dockerfiles to be built
- **api**: RESTful API to submit a dockerfile
- **dispatcher**: build the dockerfiles from the queue

Building a docker image may take a while, that's why it happens async. A queue is  used to create backpressure and decouple the _api_ and the _dispatcher.

## Setup

The following env variables need to be set to access the registry.
```
export REGISTRY_USER_ID=...
export REGISTRY_PASSWORD=...
export REGISTRY_URL=...
```
Tested with _dockerhub_ registry (https://index.docker.io/v1/).

**NOT TESTED and NOT needed for local default setup**: The documentation of the used docker client lib mentions the following env variables:
- DOCKER_HOST to set the url to the docker server
- DOCKER_API_VERSION to set the version of the API to reach, leave empty for latest
- DOCKER_CERT_PATH to load the TLS certificates from
- DOCKER_TLS_VERIFY to enable or disable TLS verification, off by default

## Starting the service

### Local

You can run the service locally with
```
go run ./image-service.go
```

### Docker

Alternatively, you can use docker, to build the image run
```
docker build -t image-service .
```

To run the _image-service_ on an unix like system and sharing the host docker socket run
```
docker run --rm -e REGISTRY_USER_ID=... -e REGISTRY_PASSWORD=... -e REGISTRY_URL=... -v "/var/run/docker.sock:/var/run/docker.sock:rw" -p 8080:8080 image-service
```

## Usage

_curl_ can be used to send a request to the api
```
curl -X POST http://localhost:8080/build \
  -F "Dockerfile=@./pkg/test-dockerfiles/alpine-dockerfile" \
  -H "Content-Type: multipart/form-data"
```

The images show up locally and in the docker registry.
