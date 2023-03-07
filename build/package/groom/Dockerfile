FROM golang:1.20.1-alpine as build
ARG script
LABEL org.opencontainers.image.authors="Florent Carré"
LABEL org.opencontainers.image.source=https://github.com/fc92/poker
LABEL org.opencontainers.image.description="poker image to choose poker room in k8s deployment context"
LABEL org.opencontainers.image.licenses=MIT

WORKDIR /go/src
COPY . /go/src
COPY build/package/groom/groom.sh /
COPY build/package/groom/clients.sh /
RUN CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@latest
# download tty2web binary
# ADD https://github.com/kost/tty2web/releases/download/v2.7.0/tty2web_linux_amd64 /tty2web
COPY tty2web /tty2web

RUN go mod download

RUN CGO_ENABLED=0 go build -gcflags="all=-N -l" -o /groom /go/src/cmd/groom/groom.go

# prepare a script /clients.sh to launch the server as well as tty2web to welcome incoming clients
# make the script and binary executable
RUN cp /go/bin/dlv / && \
    chmod +x /dlv && \
    chmod +x /groom && \
    chmod +x /groom.sh && \
    chmod +x /clients.sh && \
    chmod +x /tty2web

# avoid running as root
USER 1000:1000

WORKDIR /

# if manual debugger is needed overwrite CMD outside of Dockerfile
CMD ["/groom"]