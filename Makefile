# 프로젝트 변수 설정
PROJECT_NAME := todo-server
PKG := cmd/main.go
BINARY_NAME := $(PROJECT_NAME)

# Go 관련 변수 설정
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOCMD := go

# 빌드 작업 정의
build:
	$(GOCMD) build -o $(BINARY_NAME) $(PKG)

# 테스트 작업 정의
test:
	$(GOCMD) test -v ./...

# 실행 작업 정의
run: build
	./$(BINARY_NAME)

# 의존성 설치 작업 정의
install_deps:
	$(GOCMD) mod download

# 실행 파일 삭제 작업 정의
clean:
	rm -f $(BINARY_NAME)

# 실행 파일 설치 작업 정의
install: build
	mv $(BINARY_NAME) $(GOBIN)

# 모든 작업을 수행하는 기본 작업 정의
all: install_deps build test