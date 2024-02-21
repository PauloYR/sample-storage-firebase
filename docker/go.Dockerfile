FROM golang:alpine3.18 AS build-stage

WORKDIR /app/src/sample-storage

ENV GOPATH=/app

COPY . .

RUN chmod +x /app/src/sample-storage

RUN go test ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o /sample-storage

FROM alpine:3.19.0 AS build-release-stage

ENV TZ=America/Sao_Paulo

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app/src/sample-storage

COPY /assets /app/src/sample-storage/assets

RUN chmod +x /app/src/sample-storage

COPY --from=build-stage /sample-storage /app/src/sample-storage/sample-storage

EXPOSE 8080

ENTRYPOINT ["/app/src/sample-storage/sample-storage"]