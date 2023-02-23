FROM golang:1.20.1-alpine as build
ARG script
LABEL org.opencontainers.image.authors="Florent Carré"
LABEL org.opencontainers.image.source=https://github.com/fc92/poker
LABEL org.opencontainers.image.description="poker image running client and server components"
LABEL org.opencontainers.image.licenses=MIT

# build poker binary using multi stage build for go
WORKDIR /go/src
COPY . /go/src
RUN CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@latest

RUN go mod download

RUN CGO_ENABLED=0 go build -gcflags="all=-N -l" -o /groom /go/src/cmd/groom/groom.go

# prepare a script /clients.sh to launch the server as well as tty2web to welcome incoming clients
# make the script and binary executable
RUN cp /go/bin/dlv / && \
    chmod +x /dlv && \
    chmod +x /groom

# avoid running as root -- NEED overwrite for debug
USER 1000:1000

WORKDIR /

# if manual debugger is needed overwrite CMD outside of Dockerfile
CMD ["/groom"]