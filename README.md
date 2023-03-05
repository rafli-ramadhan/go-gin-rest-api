# Golang + Gin + Gorm + Nodemon + PostgreSQL + Swagger (REST API for PT. Phincon Attendance Apps)

Requirement :
- Go 1.16
- Node.js
- PostgreSQL

# Start 🚀

## Install Modules

```bash
go mod download && go mod tidy && go mod verify
```

If the message below was shown, do the next step.
```
go: finding module for package github.com/forkyid/go-rest-api/docs
github.com/forkyid/go-rest-api/src/route imports
        github.com/forkyid/go-rest-api/docs: no matching versions for query "latest"
```

## Swagger Installation and Swag Initialization

```bash
go install github.com/swaggo/swag/cmd/swag@v1.6.7
```

```bash
swag init -g src/main.go
```

```bash
go mod tidy
```

## Install Nodemon

```bash
npm install
```

or

```bash
npm install -g nodemon
```

## Running the Server

### Go run + Nodemon

```bash
nodemon --exec go run src/main.go --signal SIGTERM
```

Swagger API Documentation URL:
```url
http://localhost:5000/swagger/index.html#/
```

![1](/images/1.png)
![2](/images/2.png)

### Docker

```bash
docker-compose up --build
```

## Repository Structure

```bash
.
├── .github
│   └── PULL_REQUEST_TEMPLATE.md
├── database-migrations
│   ├──examples
│   └──README.md
├── src
│   ├── connection
│   │   └── connection.go
│   ├── constant
│   │   └── constant.go
│   ├── controller
│   │   └── v1
│   │       ├── account
│   │       │   └── account.go
│   │       ├── auth
│   │       │   └── auth.go
│   │       └── location
│   │           └── location.go
│   ├── http
│   │   ├── account.go
│   │   ├── auth.go
│   │   └── location.go
│   ├── model
│   │   ├── account.go
│   │   └── location.go
│   ├── pkg
│   │   ├── bcrypt
│   │   │   └── bcrypt.go
│   │   └── jwt
│   │       └── jwt.go
│   ├── repository
│   │   └── v1
│   │       ├── account
│   │       │   └── account.go
│   │       └── location
│   │           └── location.go
│   ├── routes
│   │   └── main.go
│   ├── service
│   │   └── v1
│   │       ├── account
│   │       │   └── account.go
│   │       └── location
│   │           └── location.go
│   └── main.go
├── .env.example
├── .gitignore
├── go.mod
├── LICENSE
└── README.md
```
