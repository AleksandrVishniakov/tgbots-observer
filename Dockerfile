# Stage 1: building
FROM golang:1.24-alpine AS builder
WORKDIR /build
RUN apk update && apk add --no-cache git
COPY go.mod  go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app cmd/main.go

# Stage 2: starting
FROM alpine
COPY --from=builder app /bin/app
ENTRYPOINT ["/bin/app"]
