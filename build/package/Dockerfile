FROM ubuntu:xenial

RUN echo "deb [arch=amd64] http://repo.sawtooth.me/ubuntu/ci xenial universe" >> /etc/apt/sources.list \
 && echo "deb http://archive.ubuntu.com/ubuntu xenial-backports universe" >> /etc/apt/sources.list \
 && apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys 8AA7AF1F1091A5FD \
 && apt-get update \
 && apt-get install -y -q --allow-downgrades \
    wget \
    build-essential \
    git \
    libssl-dev \
    libzmq3-dev \
    openssl \
    python3-grpcio-tools=1.1.3-1 \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/* \
 && wget -c https://dl.google.com/go/go1.11.5.linux-amd64.tar.gz \
 && tar -C /usr/local -xzf go1.11.5.linux-amd64.tar.gz

ENV GOPATH=/go
ENV PATH=$PATH:/usr/local/go/bin

RUN mkdir -p $GOPATH/src/github.com/BadgeForce/credential-template-engine
WORKDIR $GOPATH/src/github.com/BadgeForce/credential-template-engine
COPY . .

RUN make deps
RUN make build

EXPOSE 4004/tcp
ENTRYPOINT ["./credential-template-engine"]
