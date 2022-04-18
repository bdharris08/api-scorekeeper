# syntax=docker/dockerfile:1

FROM golang:alpine as builder

# Install git. (alpine image does not have git in it)
RUN apk update && apk add --no-cache git

WORKDIR /app

# Download dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/main .

# Build a minimal image with just the binary
FROM scratch

COPY --from=builder /app/bin/main .

CMD ["./main"]