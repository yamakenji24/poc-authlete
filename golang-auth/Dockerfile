FROM golang:1.24.1-alpine

WORKDIR /app

# Goのモジュール設定
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o main ./main.go

EXPOSE 3000

CMD ["./main"]