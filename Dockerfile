FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY main.go .
RUN go build -o arashpanel .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/arashpanel .
EXPOSE 8080
CMD ["./arashpanel"]
