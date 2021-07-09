FROM golang:1.16-alpine as builder
RUN mkdir -p /usr/local/src
COPY . /usr/local/src
WORKDIR /usr/local/src/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /discovery-test .

FROM scratch
COPY --from=builder /discovery-test .
ENTRYPOINT ["/discovery-test"]