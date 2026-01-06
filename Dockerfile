# build
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /forum

# dependencies first (cache-friendly)
COPY go.mod go.sum ./
RUN go mod download

# source code
COPY . .

# build binary
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -o forum-server ./cmd/web/main.go

#runtime

FROM alpine:3.20

RUN apk add --no-cache ca-certificates sqlite \
    && mkdir -p /forum/data \
    && mkdir -p /forum/ui/static/uploads

WORKDIR /forum

# binary
ENV DB_PATH=/forum/data/forum.db

COPY --from=builder /forum/forum-server .


COPY data/forum.db /forum/data/forum.db

COPY internal/db/migrations ./internal/db/migrations
COPY ui ./ui

EXPOSE 8080

CMD ["./forum-server"]
