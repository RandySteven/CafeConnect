FROM golang:1.22.0 as builder

WORKDIR /app

COPY . .

RUN go mod tidy && go mod vendor

RUN go build -o /bin/http cmd/http -config files/yml/cofeConnect.stg.yaml

