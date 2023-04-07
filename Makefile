install:
	@go get github.com/spf13/viper
	@go get gorm.io/gorm
	@go get gorm.io/driver/postgres
	@go get github.com/gin-gonic/gin

build:
	@go build -o bin/app -v ./cmd/main.go


run:
	@go run ./cmd/main.go