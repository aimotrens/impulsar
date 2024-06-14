FROM golang:1.22@sha256:4950c1cce16bb754e23ee70d25a67e906149d0faedc0aaeed49c075b696fa889 AS builder

ARG IMPULSAR_VERSION

WORKDIR /src
ADD go.mod go.sum ./
RUN go mod download

ADD . .
RUN go build -ldflags "-X \"main.impulsarVersion=${IMPULSAR_VERSION}\" -X \"main.compileDate=$(date)\"" -o impulsar

# ---

FROM alpine:3.20@sha256:77726ef6b57ddf65bb551896826ec38bc3e53f75cdde31354fbffb4f25238ebd

COPY --from=builder /src/impulsar /usr/local/bin/impulsar

ENTRYPOINT ["/usr/local/bin/impulsar"]
