FROM golang:latest as build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download -x
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/main ./cmd/main.go
 
FROM alpine:latest as run
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=build /app/main /app/main
COPY --from=build /go/bin/migrate /usr/local/bin/migrate

EXPOSE 8080
CMD ["/app/main"]