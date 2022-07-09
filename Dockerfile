FROM golang:1.18.3 AS builder
WORKDIR /app
COPY . .
RUN go build -o q-server .

FROM ubuntu:22.04
RUN apt-get update && apt-get install -y ca-certificates
COPY --from=builder /app/q-server /app/q-server
WORKDIR /app
ENTRYPOINT [ "q-server" ]