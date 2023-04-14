FROM golang:1.20-alpine AS build

ARG GH_TOKEN=""

RUN apk add build-base git linux-headers libc-dev

WORKDIR /work
COPY go.mod go.sum /work/
COPY errors/go.mod errors/go.sum /work/errors/
COPY math/go.mod math/go.sum /work/math/
COPY api/go.mod api/go.sum /work/api/
COPY core/go.mod core/go.sum /work/core/
COPY depinject/go.mod depinject/go.sum /work/depinject/
COPY tools/rosetta/go.mod tools/rosetta/go.sum /work/tools/rosetta/
COPY collections/go.mod collections/go.sum /work/collections/
COPY store/go.mod store/go.sum /work/store/
COPY log/go.mod log/go.sum /work/log/
COPY x/tx/go.mod x/tx/go.sum /work/x/tx/
COPY x/upgrade/go.mod x/upgrade/go.sum /work/x/upgrade/

RUN go env -w GOPRIVATE="github.com/bnb-chain/*"
RUN git config --global url."https://${GH_TOKEN}@github.com".insteadOf "https://github.com"

RUN go mod download

COPY ./ /work
RUN LEDGER_ENABLED=false make clean build


FROM alpine AS run
RUN apk add bash curl jq libstdc++
EXPOSE 26656 26657
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start", "--log_format", "plain"]
STOPSIGNAL SIGTERM
VOLUME /simd
WORKDIR /simd

ENV CGO_CFLAGS="-O -D__BLST_PORTABLE__"
ENV CGO_CFLAGS_ALLOW="-O -D__BLST_PORTABLE__"

COPY contrib/images/simd-env/wrapper.sh /usr/bin/wrapper.sh
COPY --from=build /work/build/simd /simd/
