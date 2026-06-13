FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.work go.work.sum ./
COPY go.mod go.sum ./
COPY generated/go ./generated/go
COPY service/order ./services/order

workdir /src/services/order
RUN go build -o /out/order ./cmd/ordder

FROM alpine:3.22

COPY --from=build /out/order /usr/local/bin/order

EXPOSE 8081
