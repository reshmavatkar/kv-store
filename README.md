# Simple Decomposed Key-Value Store in Go

This project implements a simple decomposed Key-Value store in Go with two services communicating over gRPC:

- **REST API Service:** Exposes a JSON REST API as the primary public interface.
- **gRPC Store Service:** Implements the Key-Value store logic with Set, Get, and Delete operations.

Both services run independently in Docker containers.

## Project Structure

├── api/ # REST API service

├── store/ # gRPC Key-Value store service

├── protos/ # Protobuf definitions

├── test/ # Integration tests

├── docker-compose.yml

└── README.md


## Prerequisites
- Docker and Docker Compose installed.
- Go environment set up
- Protoc installed and go plugin for protoc https://grpc.io/docs/languages/go/quickstart/

## Features

- `PUT /store` — Store a string value at a given key
- `GET /store/:key` — Retrieve the value for a key
- `DELETE /store/:key` — Delete a key-value pair
- Internal communication between services via gRPC
- In-memory map with read/write mutex for thread safety


## Getting Started

###  Clone the repository

```
$ git clone https://github.com/reshmavatkar/kv-store.git
$ cd kv-store
```
### Build and Run

Build and start both services
```
docker-compose up --build -d
```
This will:

- Build both containers (grpc-store and rest-api)
- Start grpc-store on port 50051
- Start rest-api on http://localhost:8080/, which talks to the gRPC store internally

Verify containers:
```
docker ps
```

### API Endpoints
Store a key-value pair.

```
PUT /store

Request Body:

{
  "key": "greeting",
  "value": "Hello World"
}

Response:
{
    "status":"ok"
}
```

Fetch a value for a key.
```
GET /store/:key

Response:

{
  "value": "Hello World",
  "message": "OK"
}
```

Delete a key from the store.
```
DELETE /store/:key

Response:

{
    "status":"deleted"
}
```

### Test REST API manually with curl:
Put key value
```
$ curl -X PUT http://localhost:8080/store \
  -H "Content-Type: application/json" \
  -d '{"key": "foo", "value": "bar"}'
```

Get a value for a key 
```
$ curl http://localhost:8080/store/foo
```

Delete a key
```
$ curl -X DELETE http://localhost:8080/store/foo
```

### Clean up
```
$ docker-compose down --rmi all --volumes --remove-orphans
$ go clean ./...
```

## How to use make
- `make build` to build Docker images for REST API and gRPC store.
- `make up` for start services in detached mode.
- `make down` to stop and remove containers, networks.
- `make test` to run all integration tests.
- `make clean` to clean build artifacts.

