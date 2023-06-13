# go-erply-proxy

This repository contains proposed solution for a task received from Erply.

## Task

### Read/Write customer data over API

The application uses `github.com/erply/api-go-wrapper` to read and write brands.

### Simple authentication of requests (token)

It is needed to add http header authorization with the value.

```bash
# Authorization: <sessionKey>:<client code>

curl -X 'GET' \
  'http://localhost:8081/v1/brand?skip=0&take=28' \
  -H 'accept: application/json' \
  -H 'Authorization: f3d39852878ffab3836d6e0da80c26629926eeab3e4:000000'
```

### Use database for local storage (cache)

This project utilizes Redis as a local storage solution. Data is stored using composite keys, which consist of the URI path and URI parameters. For example, the URI path could be /v1/brand, and the associated parameters may include skip=0_take=25.

During a GET request, the system first checks the cache for the requested data. If no results are found in the cache, the data is retrieved from the Erply API and directly served to the API consumer. Simultaneously, the response is saved to the cache in the background.

When a new record is added for the /v1/brand URI, the cache is cleared for the brands category key. This ensures that only the relevant brand data is removed from the cache.

To ensure data freshness, the entire cache is flushed every 30 seconds, keeping the data up to date.

This caching mechanism optimizes data retrieval by minimizing reliance on the Erply API and significantly reducing response times for subsequent requests.

### Includes simple documentation (Readme.md)
Check.

### Includes swagger or apidoc for functions

SwaggerUI is served from the server root. Additionally, a Postman collection can be found in the [./api](./api) directory.

### Includes Unit tests

Unit tests have been implemented and are passing for successful create and read list requests.

## Run
```bash
docker compose build

docker compose up               # Run all, press Ctrl+C to quite

docker compose down             # Tear down network
```

Open http://localhost:8081/ to explore the aPI using the swagger UI.


## Develop

```bash

docker compose -f docker-compose-redis.yaml up  # Run Redis

go get -d -v ./... && go install -v ./...       # Prep Go

go test -v ./...                                # Run unit tests

go run ./main.go                                # Run Go App
```

## Known issues

Pagination doesn't seem to work correctly. Attempts were made with recordsOnPage and pageNo, but without success. Further investigation is required.
