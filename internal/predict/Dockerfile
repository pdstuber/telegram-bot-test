FROM golang as BUILDER
ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64
ENV TENSORFLOW_VERSION 1.14.0
WORKDIR /usr/src/app
RUN curl -L \
  "https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-cpu-linux-x86_64-${TENSORFLOW_VERSION}.tar.gz" \
  | tar -C "/usr/local" -xz && ldconfig
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o isit-a-cat-predict-go .

FROM debian:buster-slim
RUN apt-get update && apt-get install libc6-dev ca-certificates -y
RUN update-ca-certificates
COPY --from=BUILDER /usr/local/include/tensorflow /usr/local/include/tensorflow
COPY --from=BUILDER /usr/local/lib/libtensorflow* /usr/local/lib/
RUN ldconfig
COPY --from=BUILDER /usr/src/app/isit-a-cat-predict-go .
ENTRYPOINT ["./isit-a-cat-predict-go"]
