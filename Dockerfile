
FROM golang:latest as builder

WORKDIR /build
COPY main.go /build/
COPY go.sum /build/
COPY go.mod /build/

RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix ./main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /build/pubsub-to-pubsub /app/pubsub-to-pubsub
CMD ./pubsub-to-pubsub
