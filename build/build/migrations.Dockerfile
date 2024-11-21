FROM alpine:latest

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cache/apk/*
ADD https://github.com/pressly/goose/releases/download/v3.22.1/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose

COPY ./migrations ./migrations
CMD ["goose", "-dir", "/migrations", "postgres", "postgres://user:password@postgres:5432/songs?sslmode=disable", "up" ]
