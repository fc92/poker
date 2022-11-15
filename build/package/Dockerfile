FROM golang:1.19.3 as build
LABEL org.opencontainers.image.authors="Florent Carré"
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

FROM gcr.io/distroless/static-debian11

COPY --from=build /go/bin /
CMD ["/poker"]