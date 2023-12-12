FROM golang:1.21 AS builder

WORKDIR /src
ADD go.mod go.sum ./
RUN go mod download

ADD . .
RUN go build -o impulsar

# ---

FROM alpine:3.19

COPY --from=builder /src/impulsar /usr/local/bin/impulsar

ENTRYPOINT ["/usr/local/bin/impulsar"]
