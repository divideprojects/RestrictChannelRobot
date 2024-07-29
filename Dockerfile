# Build Stage: Build bot using the alpine image, also install doppler in it
FROM golang:1.22-alpine AS builder
RUN apk add --update --no-cache git upx
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -o out/RestrictChannelRobot -ldflags="-w -s" .
RUN upx --brute out/RestrictChannelRobot

# Run Stage: Run bot using the bot and doppler binary copied from build stage
FROM alpine:3.20.2
COPY --from=builder /app/out/RestrictChannelRobot /
CMD ["/RestrictChannelRobot"]

LABEL org.opencontainers.image.authors="Divanshu Chauhan <divkix@divkix.me>"
LABEL org.opencontainers.image.url="https://divkix.me"
LABEL org.opencontainers.image.source="https://github.com/divideprojects/RestrictChannelRobot"
LABEL org.opencontainers.image.title="Restrict Channel Robot"
LABEL org.opencontainers.image.description="Official Restrict Channel Bot Docker Image"
LABEL org.opencontainers.image.vendor="Divkix"
