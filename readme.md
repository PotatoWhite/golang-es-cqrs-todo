# Golang - RESTFul API

## Description

- Todo List를 관리하는 RESTFul API

## pre-requisites

### Required

- Go 1.20.1
- Postgres
- Docker
- Docker Compose

### Optional

- Restclient

1## Installation

### Go module : Dependency Management

```bash
go mod init github.com/potato/simple-restful-api
```

### Viper : Configuration Management

```bash
go get github.com/spf13/viper
```

### GORM : ORM for Golang

```bash
go get gorm.io/gorm
```

### GORM Postgres Driver : Postgres Driver for GORM

```bash
go get gorm.io/driver/postgres
```

#### docker run postgres

```bash
docker run -p 5432:5432 -e POSTGRES_PASSWORD=1234 -e POSTGRES_USER=potato  -e POSTGRES_DB=golang --name psql_golang_study -d postgres
```

#### create account and database

```sql
CREATE USER todo_account PASSWORD '1234' SUPERUSER;
CREATE DATABASE todo_db OWNER todo_account;
CREATE USER token_account PASSWORD '1234' SUPERUSER;
CREATE DATABASE token_db OWNER token_account;
```

### Gin-Gonic : Web Framework

```bash
go get github.com/gin-gonic/gin
```

### UUID : UUID generator

```bash
go get github.com/google/uuid
```

### sarama : Kafka Client

```bash
go get github.com/Shopify/sarama
```

### mongodb

```bash
go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/bson
```

### mongo-db docker

```bash
docker run -d --name mongodb \
  -e MONGO_INITDB_ROOT_USERNAME=root \
  -e MONGO_INITDB_ROOT_PASSWORD=1234 \
  -e MONGO_INITDB_DATABASE=todo \
  -e MONGO_NON_ROOT_USERNAME=todo_account \
  -e MONGO_NON_ROOT_PASSWORD=1234 \
  -p 27017:27017 \
  --name mongo_golang_study \
  mongo
```

## Package Structure
