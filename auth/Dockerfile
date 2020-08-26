FROM golang:1.14.4 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o auth-svc .

FROM scratch

COPY --from=builder /build/auth-svc .
COPY --from=builder /build/config ./config
COPY --from=builder /build/docs ./docs

EXPOSE 8003
ENTRYPOINT ["./auth-svc"]
