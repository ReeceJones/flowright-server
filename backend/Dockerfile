FROM golang:alpine

RUN apk add build-base

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /usr/local/bin/ ./...

EXPOSE 8090

ENV PATH $PATH:/usr/local/bin/flowright

CMD ["backend", "serve", "--http", "0.0.0.0:8090"]