help: ## 도움말
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
.PHONY: help

run: swag ## 로컬 실행
	@echo "run the main.go file"
	@go mod tidy && go mod download && go run ./cmd/app/main.go
.PHONY: run

swag: ### 스웨거 초기화
	@echo "swag init"
	@swag init -g cmd/app/main.go
.PHONY: swag

mock: ### 목커리 실행
	@mockery
.PHONY: mock

test: ### 모든 테스트 실행
	go test -v -cover -race ./internal/... ./domain/...
.PHONY: test

docker-build: swag ### 도커 컴포즈 빌드
	docker build -t classting-app .

compose-clean: compose-down ### app 컨테이너 삭제
	docker rmi --force classting-app:latest

compose-up: docker-build ### app을 빌드하고 mysql 5.7과 실행
	docker-compose up -d
.PHONY: compose-up

compose-down: ### app 실행 종료
	docker-compose down
.PHONY: compose-down