# build stage
FROM golang:1.21-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org
RUN go mod download
COPY . .
WORKDIR /src/cmd/server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/server ./

# final
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/server /usr/local/bin/server
EXPOSE 8081
ENV PORT=8081
ENTRYPOINT ["/usr/local/bin/server"]
