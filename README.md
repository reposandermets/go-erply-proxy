# go-erply-proxy

## Authorization

It is needed to add http header authorization with the value

```
<sessionKey>:<client code>
```


## Run
```bash
docker compose build

docker compose up               # Run all, press Ctrl+C to quite

docker compose down             # Tear down network
```

## Develop

```bash

docker compose -f docker-compose-redis.yaml up  # Run Redis

go get -d -v ./... && go install -v ./...       # Prep Go

go test -v ./...                                # Run tests

go run ./main.go                                # Run Go App
```
