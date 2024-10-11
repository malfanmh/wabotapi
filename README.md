# wabotapi

This repository contains a Golang application that can be deployed using Docker or on a Linux environment without containers.

## Features

- Built using Go
- Dockerized for container-based deployment
- Supports deployment on a Linux environment
- Uses `tzdata` for timezone support

## Requirements

- Go 1.21+
- Docker (if using container-based deployment)
- MySQL (for database dependencies)

## Project Structure

- `cmd/main.go`: Entry point for the application
- `go.mod` and `go.sum`: Go module dependencies
- `Dockerfile`: Instructions to build and run the containerized application

---

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/your-repo/wabotapi.git
cd wabotapi
````

### Configuration
Ensure your environment variables are correctly set:

```bash
export APP_PORT=8080
export MYSQL_DSN="user:password@tcp(localhost:3306)/dbname"
export WA_BUSINESS_ACCOUNT_ID="your_wa_business_account_id"
export WA_ACCESS_TOKEN="your_wa_access_token"
export WA_BASE_URL="https://api.whatsapp.com"
export WA_SECRET="your_wa_secret"
export FINPAY_BASE_URL="https://finpay.id"

```
---
# Deployment
### A. Docker-Based Deployment
##### 1. Build Docker Image
   To build the Docker image, run:

```bash
docker build -t wabotapi .
```

##### 2. Run Docker Container
   Run the application in a container using:

```bash
docker run -d -p 8080:8080 --name my-golang-app -e APP_PORT=8080 -e MYSQL_DSN="your_mysql_dsn" my-golang-app
```

The app will be accessible at http://localhost:8080.


##### 3. Using Docker Compose (Optional)
   You can also run the app using Docker Compose if you have multiple services like MySQL. Create a docker-compose.yml file and start it with:

```bash
docker-compose up --build
```

### B. Linux Environment Deployment (Without Docker)
##### 1. Install Go and Dependencies
   Ensure Go is installed on your Linux machine. If not, install it by following this guide.

Install necessary dependencies:

```bash
sudo apt update
sudo apt install -y tzdata
```
##### 2. Build the Application
```bash
go mod download
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main ./cmd/main.go
```
##### 3. Run the Application
   Once the binary is built, run the application:

```bash
./main
```
The application will start, and you can access it at http://localhost:8080.

---

## Database Migrations
This application uses golang-migrate for managing database migrations. To run migrations:

#### 1. Install migrate if running without Docker:
```bash
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
#### 2. Run migrations:
```bash
migrate -path ./migrations -database "mysql://user:password@tcp(localhost:3306)/dbname" up
```


### ERD

```mermaid
classDiagram
    direction BT

    class clients {
        +varchar(100) name
        +varchar(100) finpay_merchant_id
        +varchar(255) finpay_secret
        +varchar(255) finpay_callback_url
        +varchar(30) wa_phone
        +varchar(100) wa_phone_id
        +varchar(100) wa_business_account_id
        +timestamp created_at
        +timestamp updated_at
        +bigint id
    }

    class customers {
        +bigint client_id
        +varchar(30) wa_id
        +varchar(100) email
        +varchar(100) full_name
        +varchar(1000) address
        +date birth_date
        +varchar(100) identity_number
        +varchar(20) identity_type  /* oneof[ktp,kitas,sim] */
        +varchar(1) gender
        +int access  /* oneof [0: public, 1:registered, 2:activated] */
        +timestamp created_at
        +timestamp updated_at
        +bigint id
    }

    class message_actions {
        +bigint message_id
        +varchar(100) slug
        +varchar(20) title
        +varchar(255) description
        +tinyint display
        +int access  /* oneof [0:public, 1:registered, 2:activated] */
        +varchar(20) seq
        +timestamp created_at
        +timestamp updated_at
        +bigint id
    }

    class message_flows {
        +bigint client_id
        +varchar(100) keyword
        +bigint message_id
        +int access
        +tinyint display
        +varchar(100) seq
        +timestamp created_at
        +timestamp updated_at
        +tinyint validate_input
        +tinyint checkout
        +bigint id
    }

    class messages {
        +bigint client_id
        +varchar(100) slug
        +varchar(100) type  /* oneof [text, button, interactive] */
        +int access  /* oneof [0:public, 1:registered, 2:activated] */
        +varchar(50) button
        +text header_text
        +tinyint preview_url
        +text body_text
        +text footer_text
        +tinyint with_metadata
        +timestamp created_at
        +timestamp updated_at
        +bigint id
    }

    class payments {
        +timestamp created_at
        +timestamp updated_at
        +timestamp expired_at
        +bigint customer_id
        +varchar(50) ref_id
        +varchar(50) payment_provider
        +varchar(50) payment_ref_id
        +varchar(20) payment_type  /* oneof[va, bank_transfer, qris, indomaret] */
        +varchar(100) payment_code
        +varchar(100) payment_item
        +varchar(20) status  /* oneof [PENDING, PAID, EXPIRED] */
        +decimal(30,2) amount
        +decimal(30,2) fee
        +bigint client_id
        +bigint id
    }

    class products {
        +timestamp created_at
        +timestamp updated_at
        +bigint client_id
        +bigint customer_id
        +varchar(200) name
        +text description
        +varchar(200) slug
        +decimal(30,2) price
        +decimal(30,2) stock
        +varchar(250) image
        +bigint id
    }

    class session {
        +bigint client_id
        +varchar(30) wa_id
        +int access
        +varchar(10) seq
        +varchar(100) slug
        +varchar(100) input
        +timestamp created_at
        +timestamp updated_at
        +timestamp expired_at
        +bigint id
    }

    customers --> clients
    message_actions --> messages
    message_flows --> clients
    message_flows --> messages
    messages --> clients
    payments --> clients
    payments --> customers
    products --> clients
    products --> customers
    session --> clients

```
