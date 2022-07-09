# build binary
FROM golang:1.18.3-alpine AS builder

RUN apk add --no-cache gcc musl-dev linux-headers

WORKDIR /app

COPY . .

RUN go build -o q-server .

# copy binary to main container
FROM alpine:3.16.0

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/q-server /app/q-server

WORKDIR /app

ENTRYPOINT [ "/app/q-server" ]