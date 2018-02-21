FROM gliderlabs/alpine:3.4

env GOROOT /go
ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /go/lib/time/zoneinfo.zip
RUN apk --no-cache add ca-certificates openssl && \
    wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://raw.githubusercontent.com/sgerrand/alpine-pkg-glibc/master/sgerrand.rsa.pub && \
    wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.25-r0/glibc-2.25-r0.apk && \
    apk --no-cache add glibc-2.25-r0.apk
RUN apk update \
    && apk add sqlite \
    && apk add socat
WORKDIR /shima-chain
COPY ./bin ./bin

RUN mkdir /shima-chain/db
RUN touch /shima-chain/db/sqlite3.db
RUN chmod 777 /shima-chain/
RUN chmod 777 /shima-chain/db/
RUN chmod 777 /shima-chain/db/sqlite3.db

ADD migrate/ /migrate
RUN sqlite3 /shima-chain/db/sqlite3.db < /migrate/account.sql
RUN sqlite3 /shima-chain/db/sqlite3.db < /migrate/block.sql
RUN sqlite3 /shima-chain/db/sqlite3.db < /migrate/health.sql
RUN sqlite3 /shima-chain/db/sqlite3.db < /migrate/tx.sql
RUN sqlite3 /shima-chain/db/sqlite3.db < /migrate/dummy.sql

CMD /shima-chain/bin/shima-chain
