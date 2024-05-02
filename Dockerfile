FROM alpine as base
#enable the testing repository
RUN echo "@testing http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories
RUN apk add --no-cache ca-certificates postgresql-client curl tini bash gnupg taskd taskd-pki git task py3-pip

# Install virtual environment package
RUN apk add --no-cache gcc g++ build-base cmake libc-dev libffi-dev openssl-dev python3-dev py3-virtualenv

# Create a virtual environment
RUN python3 -m venv /venv

# Activate the virtual environment
ENV PATH="/venv/bin:$PATH"
RUN pip3 install timewsync

# Install timewarrior
RUN curl -L -O https://github.com/GothenburgBitFactory/timewarrior/releases/download/v1.7.1/timew-1.7.1.tar.gz
RUN tar -xvf timew-1.7.1.tar.gz
WORKDIR /timew-1.7.1
RUN cmake .
RUN make
RUN make install

FROM golang:alpine as builder
ENV GO111MODULE=on

RUN mkdir -p /app/
WORKDIR /app
ADD go.mod .
ADD go.sum .
RUN go mod download
RUN apk add --no-cache --update go gcc g++ taskd taskd-pki

ADD . .

RUN CGO_ENABLED=1 go build -o app cmd/server/main.go
RUN CGO_ENABLED=1 go build -o db cmd/db/main.go

FROM golang:alpine as timewsync
ENV GO111MODULE=on

RUN mkdir -p /app/
WORKDIR /app
RUN apk add --no-cache --update git go gcc g++ openssh

RUN git clone https://github.com/timewarrior-synchronize/timew-sync-server.git
WORKDIR /app/timew-sync-server

RUN go mod download
RUN go build -o timew-server

FROM base as prod
ENTRYPOINT ["/sbin/tini", "--"]
ENV GO_ENV=production
EXPOSE 3000

# Define volumes
VOLUME ["/app/authorized_keys", "/app/taskd/", "/app/data"]

WORKDIR /app
COPY --from=builder /app/app /app/db ./
COPY --from=timewsync /app/timew-sync-server/timew-server ./timew-server
RUN touch .env
RUN mkdir -p ./authorized_keys

ENV TASKDDATA=/app/taskd
ENV TASKD_SERVER=0.0.0.0:53589
ENV TIMEW_SERVER=http://localhost:8080
ENV TIMEW_SYNC=/app/
ENV CERT_ORGANIZATION = Flow
ENV CERT_CN = benciks.me
ENV CERT_COUNTRY = CZ
ENV CERT_STATE = Jihomoravsky
ENV CERT_LOCALITY = Brno

# Init taskd server
COPY --from=builder /app/setup/taskd.sh ./
RUN ls -la /app
RUN chmod +x ./taskd.sh

EXPOSE 53589
EXPOSE 3000
EXPOSE 8080
# Create a startup script
RUN echo "#!/bin/sh" >> /app/start.sh && \
    echo "./taskd.sh" >> /app/start.sh && \
    echo "taskdctl start" >> /app/start.sh && \
    echo "/app/timew-server start &" >> /app/start.sh && \
    echo "/app/db" >> /app/start.sh && \
    echo "/app/app" >> /app/start.sh && \
    chmod +x /app/start.sh

CMD ["/app/start.sh"]