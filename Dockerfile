FROM golang:1.21.8-alpine AS builder

COPY . /github.com/semho/chat-microservices/auth/
WORKDIR /github.com/semho/chat-microservices/auth/

RUN go mod download
RUN go build -o ./bin/auth_server cmd/server/main.go

FROM alpine:3.19.1

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cache/apk/* \

WORKDIR /root/
COPY --from=builder /github.com/semho/chat-microservices/auth/bin/auth_server .
COPY --from=builder /github.com/semho/chat-microservices/auth/entrypoint.sh .
COPY --from=builder /github.com/semho/chat-microservices/auth/migrations ./migrations

ADD https://github.com/pressly/goose/releases/download/v3.14.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose
RUN chmod +x entrypoint.sh