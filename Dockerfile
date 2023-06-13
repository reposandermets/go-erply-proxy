FROM golang:1.17-alpine AS Builder
RUN apk add --no-cache git
RUN apk add build-base
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o /go/bin/app

# Final stage
FROM alpine:latest AS Application
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache tzdata
ENV TZ Europe/Tallinn
COPY --from=builder /go/bin/app /app

ENTRYPOINT ./app
LABEL Name=ErplyProxy Version=0.0.1
EXPOSE 8081
