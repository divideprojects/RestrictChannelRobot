FROM alpine:3.18.0
COPY restrictchannelrobot /
CMD ["/restrictchannelrobot"]
LABEL org.opencontainers.image.authors="Divanshu Chauhan <divkix@divkix.me>"
LABEL org.opencontainers.image.url="https://divkix.me"
LABEL org.opencontainers.image.source="https://github.com/divkix/RestrictChannelRobot"
LABEL org.opencontainers.image.title="Restrict Channel Robot"
LABEL org.opencontainers.image.description="Official Restrict Channel Bot Docker Image"
LABEL org.opencontainers.image.vendor="Divkix"
