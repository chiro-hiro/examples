# Build Geth in a stock Go builder container
FROM golang:1.16-alpine as builder

ENV GETH_TAG="v1.10.2" GETH_NETWORK="mainnet"

RUN apk add --no-cache make gcc musl-dev linux-headers git

RUN cd / && \
  git clone https://github.com/ethereum/go-ethereum && \
  cd /go-ethereum/ && \
  git reset --hard ${GETH_TAG} && \
  make geth

# Pull Geth into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go-ethereum/build/bin/geth /usr/local/bin/

EXPOSE 8546 30303 30303/udp

CMD geth --${GETH_NETWORK} --http --http.addr "0.0.0.0" --http.port 8546 --http.corsdomain "*" --syncmode light --light.maxpeers 255 --maxpeers 255
