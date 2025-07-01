FROM golang:1.23.0-alpine

WORKDIR /app

COPY ../../go.mod go.sum ./
RUN go mod tidy && go mod vendor

COPY ../.. .

COPY ../yml/cofeConnect.docker.yaml /files/yml/cofeConnect.docker.yaml

RUN go build -o /bin/consumer ./cmd/main/consumer

CMD ["/bin/consumer", "-config", "/files/yml/cofeConnect.docker.yaml"]
