# golang:1.21-alpine3.18
FROM --platform=${TARGETPLATFORM:-linux/amd64} golang:1.21.4-bookworm as builder

ENV USER=appuser
ENV UID=10001

ARG TARGETARCH TARGETOS

ENV LIBTENSORFLOW_FILENAME="libtensorflow-2.14.1-${TARGETARCH}.tar.gz"

WORKDIR /app

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "$(pwd)" \
    --no-create-home \
    --uid "$UID" \
    "$USER"

ADD build/${LIBTENSORFLOW_FILENAME} /usr/local/

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod verify

RUN ldconfig /usr/local/lib

ADD . .

RUN --mount=type=cache,target="/root/.cache/go-build" \
    GOOS=$TARGETOS GOARCH=$TARGETARCH GOCACHE=/root/.cache/go-build go build -o /go/bin/telegram-bot-test

FROM --platform=${TARGETPLATFORM:-linux/amd64} debian:bookworm-slim

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /usr/local/lib/libtensorflow* /usr/local/lib/
COPY --from=builder /usr/local/include/tensorflow /usr/local/include/

RUN ldconfig /usr/local/lib

COPY --from=builder /go/bin/telegram-bot-test /go/bin/telegram-bot-test

RUN mkdir /model
ADD build/model/* /model/

USER appuser:appuser

ENTRYPOINT ["/go/bin/telegram-bot-test"]