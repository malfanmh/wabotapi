FROM golang:latest as build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download -x
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/main ./cmd/main.go
 
FROM alpine:latest as run

WORKDIR /app

COPY --from=build /app/main /app/main
EXPOSE 8080
ENV APP_PORT=8080 \
    MYSQL_DSN=wabot_user:password@tcp(127.0.0.1:3306)/wabot_db?parseTime=true \
    WA_BUSINESS_ACCOUNT_ID= \
    WA_ACCESS_TOKEN= \
    WA_BASE_URL=https://graph.facebook.com/v20.0 \
    WA_SECRET= \
    FINPAY_BASE_URL=

CMD ["/app/main"]