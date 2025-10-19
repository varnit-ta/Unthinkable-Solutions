# build stage
FROM golang:1.24.2-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
# ensure git is available for modules using git+ URLs
RUN apk add --no-cache git && \
	go env -w GOPROXY=https://proxy.golang.org
RUN go mod download
COPY . .
WORKDIR /src/cmd
RUN mkdir -p /app && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -trimpath -ldflags="-s -w" -o /app/server .

# final
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/server /usr/local/bin/server
RUN chmod +x /usr/local/bin/server
EXPOSE 8081
ENV PORT=8081
ENTRYPOINT ["/usr/local/bin/server"]
