# Dekamond Task

A simple user registration and authentication API using Go, Echo, and SQLite.

## How it works

- POST /register - Register with mobile number, receive OTP in console
- POST /login - Login with mobile number and OTP to get JWT token
- GET /users - Get paginated list of users
- GET /users/{mobile} - Get user by mobile number

## Run locally

```bash
go run cmd/main.go
```

## Run with Docker

```bash
docker build -t dekamond-task .
docker run -p 8080:8080 dekamond-task
```

## API Documentation

Visit http://localhost:8080/swagger/ to see Swagger documentation
