FROM golang:latest AS builder

WORKDIR /app/

COPY go.mod go.mod
COPY go.sum go.sum
COPY cmd cmd
COPY pkg pkg
COPY Makefile .
RUN make go-build

RUN chmod a+x bin/api-pipeline

FROM alpine
WORKDIR /app/

COPY --from=builder /app/bin /app/

ENTRYPOINT ["/app/api-pipeline"]
