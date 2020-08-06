FROM golang:1.14.4

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

WORKDIR /app

RUN cp -r /build/config .
RUN cp -r /build/docs .
RUN cp /build/auth-svc .

EXPOSE 8003
ENTRYPOINT ["/app/auth-svc"]
