# Create base builder image
FROM golang:1.18.8-alpine3.16 AS builder
WORKDIR /go/src/github.com/chain4travel/magellan
RUN apk add --no-cache alpine-sdk bash git make gcc musl-dev linux-headers git ca-certificates g++ libstdc++


# Build app
COPY . .
RUN if [ -d "./vendor" ];then export MOD=vendor; else export MOD=mod; fi && \
    GOOS=linux GOARCH=amd64 go build -mod=$MOD -o /opt/magelland ./cmds/magelland/*.go

# Create final image
FROM alpine:3.16 as execution
RUN apk add --no-cache libstdc++
VOLUME /var/log/magellan
WORKDIR /opt

# Copy in and wire up build artifacts
COPY --from=builder /opt/magelland /opt/magelland
COPY --from=builder /go/src/github.com/chain4travel/magellan/docker/columbus/config.json /opt/config.json
COPY --from=builder /go/src/github.com/chain4travel/magellan/services/db/migrations /opt/migrations
ENTRYPOINT ["/opt/magelland"]
