FROM golang:1.21.7 as build
LABEL org.opencontainers.image.authors="fc92"
LABEL org.opencontainers.image.source=https://github.com/fc92/poker
LABEL org.opencontainers.image.description="poker image running client and server components"
LABEL org.opencontainers.image.licenses=MIT

# build poker binary using multi stage build for go
WORKDIR /go/src/app
COPY backend/go.mod ./
COPY backend/go.sum ./

RUN go mod download
COPY backend/cmd/ cmd/
COPY backend/internal/ internal/
RUN pwd
RUN ls cmd
RUN go vet -v ./...
RUN go test -v ./...

WORKDIR /go/src/app
RUN CGO_ENABLED=0 go build -o /go/bin/poker cmd/poker.go


FROM alpine:latest
WORKDIR /
# Create a non-root user and group
RUN addgroup -g 1001 appgroup && adduser -D -u 1001 -G appgroup appuser
ENV FRONTEND_DIR="/static"
RUN mkdir -p ${FRONTEND_DIR}
COPY --from=build /go/bin /
COPY "frontend/dist" $FRONTEND_DIR

USER appuser
# Set environment variables
ENV ROOM_LIST="Fire,Water,Air,Earth"
EXPOSE 8080

# Command to run the executable
CMD ["./poker", "server", "-debug"]