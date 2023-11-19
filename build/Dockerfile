FROM golang@sha256:f475434ea2047a83e9ba02a1da8efc250fa6b2ed0e9e8e4eb8c5322ea6997795 as builder

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

ENV USER=appuser
ENV UID=10001

WORKDIR /app

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "$(pwd)" \
    --no-create-home \
    --uid "$UID" \
    "$USER"

COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o /go/bin/telegram-bot-test

FROM --platform=${TARGETPLATFORM:-linux/amd64} scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /go/bin/telegram-bot-test /go/bin/telegram-bot-test

USER appuser:appuser

ENTRYPOINT ["/go/bin/telegram-bot-test"]