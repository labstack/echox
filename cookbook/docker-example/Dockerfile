FROM golang:1.15-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod init echoserver && go get github.com/labstack/echo/v4

RUN  export CGO_ENABLED=0 && go build -o /app/main.go

FROM scratch

COPY --from=builder /app/main.go /app/main.go

EXPOSE 1323

ENTRYPOINT ["/app/main.go"]
