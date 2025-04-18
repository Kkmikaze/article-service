# Article Service

Article Service is a code base for a service. It can be used 1 code base for 1 service or 1 code base for many services. Go gRPC Microservices has implement Clean Architecture pattern its mean you can use any packages or libraries you want to use.

## Table of Contents
- [Features](#features)
- [Pre-configured](#pre-configured)
- [Getting Started](#getting-started)
- [Available Scripts](#available-scripts)
- [Architecture](#architecture)

## Features
- Thorough documentation: Written with the same care as Go docs.
- Guaranteed consistency: Opinionated linting for Go integrated into Text Editor or IDE and run against staged files on pre-commit.
- Used Clean Architecture with standard project layout in Go.
- Integrated gRPC Gateway and OpenAPI Generator.
- Use Swagger UI for Implement OpenAPI/Swagger files.


## Pre-configured
- go gRPC Gateway
- protoc-gen-openapiv2
- gorm
- postgres
- Viper (for environment)
- Cobra (for execution function)

## Getting Started
Make sure you have the following installed:
- [Go](https://go.dev/doc/install)
- [Protocol Buffers](https://grpc.io/docs/languages/go/quickstart/)
- [Go gRPC Gateway](https://grpc-ecosystem.github.io/grpc-gateway/)
- [Go Migrate](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md)

#### 1. Install required Go gRPC Gateway

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest && go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@lates
```

#### 2. Clone the repository.
```bash
git clone https://github.com/Kkmikaze/article-service.git
```

#### 3. Enter your cloned directory.
```bash
cd article-service
```

#### 4. Install dependencies.
```bash
make dependencies
```

#### 5. Generate protoc and swagger openapi file
```bash
make protoc-gen
```

#### 6. Copy .env.example into .env
```bash
cp .env.example .env
```

#### 7. Run migration
```bash
make migrate-up
```

#### 8. Run on your local.
This command is a default to run development mode and wil be listen http://localhost:40000 for rpc and http://localhost:4000 for gateway
```bash
make debug
```

## Available Scripts
In the project directory, you can run:
```bash
# 1. Running in local environment
make debug

# 2. Install dependencies
make dependencies

# 3. Generate protoc and Swagger openapi file
make protoc-gen

# 4. Clean all stubs protoc and Swagger openapi file
make clean-proto

# 5. Run migration up
make migrate-up

# 6. Run migration down
make migrate-down

# 6. Running of unit test
make run-test

```

## Architecture
```
├── api
└── cmd
    ├── server/main.go
    └── root.go
├── common
├── config
├── constants
└── internal
    └── domain
        └── {service-name}
            └── {service-version}
                ├── entity
                ├── handler
                ├── repository
                ├── schema
                ├── usecase
                └── {service-name}.go
└── pkg
    ├── gateway
    ├── interceptors
    ├── orm
    ├── rpcclient
    └── rpcserver
└── proto
    └── {service-name}
        └── {service-version}
            ├── {service}_message.proto
            └── {service}_service.proto
└── scripts
    ├── protoc-gen.sh
    └── swagger-ui-gen.sh
├── stubs
├── third_party
├── .env.example
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile
└── README.md
```
