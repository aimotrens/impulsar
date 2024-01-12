FROM golang:1.21@sha256:6fbd2d3398db924f8d708cf6e94bd3a436bb468195daa6a96e80504e0a9615f2 AS builder

ARG IMPULSAR_VERSION

WORKDIR /src
ADD go.mod go.sum ./
RUN go mod download

ADD . .
RUN go build -ldflags "-X \"main.impulsarVersion=${IMPULSAR_VERSION}\" -X \"main.compileDate=$(date)\"" -o impulsar

# ---

FROM alpine:3.19@sha256:51b67269f354137895d43f3b3d810bfacd3945438e94dc5ac55fdac340352f48

COPY --from=builder /src/impulsar /usr/local/bin/impulsar

ENTRYPOINT ["/usr/local/bin/impulsar"]
