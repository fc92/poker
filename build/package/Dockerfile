FROM golang:1.19.5 as build
LABEL org.opencontainers.image.authors="Florent Carré"
LABEL org.opencontainers.image.source=https://github.com/fc92/poker
LABEL org.opencontainers.image.description="poker image to run server or client"
LABEL org.opencontainers.image.licenses=MIT

WORKDIR /go/src/app
COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY cmd/ cmd/
COPY internal/ internal/
RUN pwd
RUN ls cmd
RUN go vet -v ./cmd
RUN go test -v ./cmd

WORKDIR /go/src/app/cmd
RUN CGO_ENABLED=0 go build -o /go/bin/poker 

# for tty2web an image with a shell is required
FROM alpine:latest
USER 1000
# required for client using infocmp
RUN apk add ncurses
COPY --from=build /go/bin /
COPY https://github.com/kost/tty2web/releases/download/v2.7.0/tty2web_linux_amd64 /tty2web
RUN chmod +x /tty2web

# for tty2web start the container with dedicated parameters instead of the default CMD 
# example: docker run -it -p 8081:8081/tcp ghcr.io/fc92/poker:main "/tty2web" "--title-format" "Poker" "--permit-arguments" "-w" "-p" "8081" "/bin/sh" "-c" "/poker client -websocket 192.168.0.1:8080" &
CMD ["/poker","server","-websocket", "0.0.0.0:8080"]