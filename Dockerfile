# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

EXPOSE 8080

WORKDIR /app

COPY go.mod ./
#COPY go.sum ./
RUN go mod download

COPY car-audio-database/dist ./car-audio-database/dist

COPY server ./server
RUN go build -o /db-server ./server/*.go
CMD [ "/db-server" ]