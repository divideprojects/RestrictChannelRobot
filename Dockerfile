# Build Stage: Build bot using the alpine image, also install doppler in it
FROM golang:1.17-alpine3.14 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -o out/RestrictChannelRobot
RUN (curl -Ls --tlsv1.2 --proto "=https" --retry 3 https://cli.doppler.com/install.sh || wget -t 3 -qO- https://cli.doppler.com/install.sh) | sh

# Run Stage: Run bot using the bot and doppler binary copied from build stage
FROM gcr.io/distroless/static
COPY --from=builder /app/out/RestrictChannelRobot /
COPY --from=builder /usr/local/bin/doppler /

# Comment/Remove the line below
CMD ["/doppler", "run", "--", "/RestrictChannelRobot"]

# uncomment the line below before use, because we use doppler for secrets managment
# CMD ["/RestrictChannelRobot"]
