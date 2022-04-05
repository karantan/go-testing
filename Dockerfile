FROM golang:alpine

RUN apk add --no-cache make curl gcc libc-dev git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV DOCKER=1
ENV DEV_INFLUX_DSN=http://influxdb:8086
CMD make integration-tests
