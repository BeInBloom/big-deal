FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.work go.work.sum ./
COPY go.mod ./
COPY generated/go ./generated/go
COPY services/order ./services/order

WORKDIR /src/services/order
RUN go build -o /out/order ./cmd/order

FROM alpine:3.22

COPY --from=build /out/order /usr/local/bin/order

EXPOSE 8081

ENTRYPOINT ["order"]
