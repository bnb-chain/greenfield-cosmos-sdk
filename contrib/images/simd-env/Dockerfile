FROM golang:1.18-alpine AS build

ARG GH_TOKEN=""

RUN apk add build-base git linux-headers libc-dev

WORKDIR /work
COPY go.mod go.sum /work/
COPY db/go.mod db/go.sum /work/db/
COPY errors/go.mod errors/go.sum /work/errors/
COPY math/go.mod math/go.sum /work/math/
COPY ./ics23/go/go.mod /work/ics23/go/go.mod

RUN go env -w GOPRIVATE="github.com/bnb-chain/*"
RUN git config --global url."https://${GH_TOKEN}@github.com".insteadOf "https://github.com"

RUN go mod download
COPY ./ /work
RUN LEDGER_ENABLED=false make clean build

FROM alpine AS run
RUN apk add bash curl jq libstdc++
COPY contrib/images/simd-env/wrapper.sh /usr/bin/wrapper.sh

ENV CGO_CFLAGS="-O -D__BLST_PORTABLE__"
ENV CGO_CFLAGS_ALLOW="-O -D__BLST_PORTABLE__"

VOLUME /simd
COPY --from=build /work/build/simd /simd/
WORKDIR /simd

EXPOSE 26656 26657
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start", "--log_format", "plain"]
STOPSIGNAL SIGTERM
