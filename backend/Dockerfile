FROM golang:1.22-alpine
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /bin/localtunes ./cmd/server
EXPOSE 8080
CMD ["/bin/localtunes"]
