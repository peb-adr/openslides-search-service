FROM golang:1.24.0-alpine as base
WORKDIR /root/openslides-search-service

RUN apk add git

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY pkg pkg

# Build service in seperate stage.
FROM base as builder
RUN go build -o openslides-search-service cmd/searchd/main.go


# Test build.
FROM base as testing

RUN apk add build-base

CMD go vet ./... && go test -test.short ./...


# Development build.
FROM base as development

RUN ["go", "install", "github.com/githubnemo/CompileDaemon@latest"]
EXPOSE 9050

WORKDIR /root
COPY entrypoint.sh ./
COPY meta/search.yml .
COPY meta/models.yml .
ENTRYPOINT ["./entrypoint.sh"]

CMD CompileDaemon -log-prefix=false -build="go build -o search-service ./openslides-search-service/cmd/searchd/main.go" -command="./search-service"


# Productive build
FROM alpine:3

LABEL org.opencontainers.image.title="OpenSlides Search Service"
LABEL org.opencontainers.image.description="The Search Service is a http endpoint where the clients can search for data within Openslides."
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source="https://github.com/OpenSlides/openslides-search-service"

COPY entrypoint.sh ./
COPY meta/search.yml .
COPY meta/models.yml .
COPY --from=builder /root/openslides-search-service/openslides-search-service .
EXPOSE 9050
ENTRYPOINT ["./entrypoint.sh"]
CMD exec ./openslides-search-service
