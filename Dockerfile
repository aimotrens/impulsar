FROM golang:1.21@sha256:672a2286da3ee7a854c3e0a56e0838918d0dbb1c18652992930293312de898a6 AS builder

WORKDIR /src
ADD go.mod go.sum ./
RUN go mod download

ADD . .
RUN go build -o impulsar

# ---

FROM alpine:3.19@sha256:51b67269f354137895d43f3b3d810bfacd3945438e94dc5ac55fdac340352f48

COPY --from=builder /src/impulsar /usr/local/bin/impulsar

ENTRYPOINT ["/usr/local/bin/impulsar"]
