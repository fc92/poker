FROM golang:1.21.0 as build
LABEL org.opencontainers.image.authors="Florent Carré"
LABEL org.opencontainers.image.source=https://github.com/fc92/poker
LABEL org.opencontainers.image.description="poker image running client and server components"
LABEL org.opencontainers.image.licenses=MIT

# build poker binary using multi stage build for go
WORKDIR /go/src/app
COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY cmd/ cmd/
COPY internal/ internal/
RUN pwd
RUN ls cmd
RUN go vet -v ./...
RUN go test -v ./...

WORKDIR /go/src/app
RUN CGO_ENABLED=0 go build -o /go/bin/poker cmd/poker.go

# for tty2web an image with a shell is required
FROM alpine:3.18.3
RUN addgroup -S myapp && adduser -S myapp -G myapp

# required for client using infocmp
COPY --from=build /go/bin /

# download tty2web binary
ADD https://github.com/kost/tty2web/releases/download/v3.0.0/tty2web_linux_amd64 /tty2web

COPY build/package/poker.sh /
COPY build/package/clients.sh /
# prepare a script /clients.sh to launch the server as well as tty2web to welcome incoming clients
# make the script and binary executable
RUN chmod +x /tty2web \
    && chmod +x /clients.sh \
    && chmod +x /poker.sh

# avoid running as ROOT
USER myapp

# start server and tty2web to welcome clients on specified port
CMD ["/clients.sh", "8081"]