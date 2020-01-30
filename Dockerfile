FROM golang:1.13.5-alpine3.10 as builder
ENV GO111MODULE=on
RUN apk add --no-cache gcc musl-dev
WORKDIR /go/src/github.com/jatgam/wishlist-api/
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux go build -a -installsuffix cgo -o /wishlist_api ./

# Deploy Container
FROM alpine:3.10
ENV APP_USER=wishlist \
    APP_NAME=wishlist_api
RUN apk add --no-cache --update bash ca-certificates && \
    # Setup a base user for running applications
    adduser -u 1337 -D -h /home/${APP_USER} ${APP_USER} && \
    update-ca-certificates && \
    mkdir -p /opt/${APP_NAME} && \
    chown -R ${APP_USER}:${APP_USER} /opt/${APP_NAME} && \
    rm -rf /usr/share/man /tmp/* /var/tmp/* /var/cache/apk/*

COPY --from=builder /wishlist_api /opt/${APP_NAME}/.

USER ${APP_USER}
WORKDIR /opt/${APP_NAME}
ENTRYPOINT ["./wishlist_api"]
