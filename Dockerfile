FROM gliderlabs/alpine:3.4

env GOROOT /go
ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /go/lib/time/zoneinfo.zip
RUN apk --no-cache add ca-certificates openssl && \
    wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://raw.githubusercontent.com/sgerrand/alpine-pkg-glibc/master/sgerrand.rsa.pub && \
    wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.25-r0/glibc-2.25-r0.apk && \
    apk --no-cache add glibc-2.25-r0.apk
WORKDIR /anzu-chain
COPY ./bin ./bin
CMD /anzu-chain/bin/anzu-chain
