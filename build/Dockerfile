FROM golang:1.26@sha256:5822931cf78fe98a97edcf73a0c54c29fa2386b99c8136468e274ae9fab8cfba AS builder

ARG IMPULSAR_VERSION

WORKDIR /src
ADD go.mod go.sum ./
RUN go mod download

ADD . .
RUN go build -ldflags "-X \"main.impulsarVersion=${IMPULSAR_VERSION}\" -X \"main.compileDate=$(date +%s)\"" -o impulsar ./cmd/impulsar

# ---

FROM alpine:3.24@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b

COPY --from=builder /src/impulsar /usr/local/bin/impulsar

ENTRYPOINT ["/usr/local/bin/impulsar"]
