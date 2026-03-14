FROM golang:1.25 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /bin/server /app/server
COPY migrations /app/migrations
COPY docs/swagger.json /app/docs/swagger.json
EXPOSE 8080
ENTRYPOINT ["/app/server"]
