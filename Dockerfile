FROM golang:1.23@sha256:4a3c2bcd243d3dbb7b15237eecb0792db3614900037998c2cd6a579c46888c1e AS builder

ARG IMPULSAR_VERSION

WORKDIR /src
ADD go.mod go.sum ./
RUN go mod download

ADD . .
RUN go build -ldflags "-X \"main.impulsarVersion=${IMPULSAR_VERSION}\" -X \"main.compileDate=$(date +%s)\"" -o impulsar

# ---

FROM alpine:3.20@sha256:e417839c51cd76908ea4ec131f132ac80a9e77d345d8a4dbc3ab2ab784f8ee6c

COPY --from=builder /src/impulsar /usr/local/bin/impulsar

ENTRYPOINT ["/usr/local/bin/impulsar"]
