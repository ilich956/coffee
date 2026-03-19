# syntax=docker/dockerfile:1

FROM golang:1.23

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /hot-coffee ./cmd

EXPOSE 8080

CMD [ "/hot-coffee" ]
