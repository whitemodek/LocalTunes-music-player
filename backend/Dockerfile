FROM golang:1.21-alpine AS builder

WORKDIR /app
    
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o localtunes ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/localtunes .
COPY --from=builder /app/uploads ./uploads
COPY --from=builder /app/static ./static

EXPOSE 8080
CMD ["./localtunes"]