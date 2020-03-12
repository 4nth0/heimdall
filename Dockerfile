FROM golang:1.13.1 as builder
ARG GIT_TAG_NAME
ARG LD_FLAGS="-s -w -X main.Version=$GIT_TAG_NAME"

WORKDIR /project
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "$LD_FLAGS" -o heimdall ./

FROM alpine:3.10.2
RUN apk --update add --no-cache ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /project/heimdall heimdall

RUN pwd

CMD ./heimdall
