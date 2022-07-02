# Build Step
FROM golang:1.18-alpine AS builder

# Dependencies
RUN apk update && apk add --no-cache git upx
# Source
WORKDIR $GOPATH/src/github.com/jelliflix/jelliflix
ENV GO111MODULE on
ENV CGO_ENABLED 0
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
COPY . .

# Build
RUN go build -a -gcflags='-N -l' -installsuffix cgo -o /tmp/jelliflix
RUN upx --best --lzma /tmp/jelliflix

# Runtime Step
FROM gcr.io/distroless/static
COPY --from=builder /tmp/jelliflix /go/bin/jelliflix
EXPOSE 8080

ENTRYPOINT ["/go/bin/jelliflix"]
