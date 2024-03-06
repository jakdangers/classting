FROM golang:1.21-alpine AS builder

# 복사 할 디렉토리 설정
WORKDIR /app

# 모듈 다운로드
COPY go.mod go.sum ./
RUN go mod download

# 모든 파일 복사
ADD . /app

RUN apk add alpine-sdk
RUN go build -tags dev -v -a -ldflags="-X 'classting/config/config.configMode=dev'" -o bin/classting cmd/app/main.go

# 다단계 빌드
FROM alpine

WORKDIR /app

# 바이너리 실행 파일 복사
COPY --from=builder /app/bin/classting /app/classting
# 환경설정 파일 복사
COPY config /app/config

EXPOSE 3000

CMD ["/app/classting"]