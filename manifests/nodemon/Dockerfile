# BUILD STAGE
FROM golang:1.18 AS builder

LABEL maintainer="gdsoumya"

ARG TARGETOS=linux
ARG TARGETARCH

ADD . /server
WORKDIR /server

ENV GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH}

RUN go env
RUN CGO_ENABLED=0 go build -o /output/server -v

## DEPLOY STAGE
FROM alpine:latest

LABEL maintainer="gdsoumya"

COPY --from=builder /output/server /

CMD ["./server"]

EXPOSE 8080