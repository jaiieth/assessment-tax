FROM golang:1.22.2-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN CGO_ENABLED=0 go build -o /ktaxes-api

FROM builder as test
RUN go test -v ./...

FROM alpine:latest

COPY --from=builder /ktaxes-api /ktaxes-api

ENTRYPOINT ["/ktaxes-api"]