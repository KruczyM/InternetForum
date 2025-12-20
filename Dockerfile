FROM golang:1.24-alpine
RUN apk add --no-cache gcc musl-dev
#gcc and musl are standard C tools needed cause sqlite3 uses C lang.
WORKDIR /app
#copy dependecies first allows docker to cache downloads for faster builds later
COPY go.mod go.sum ./
RUN go mod download
#copy the rest of code
COPY . .
RUN go build -o forum-server ./cmd/server/main.go
#inform docker the party is on port 8080
EXPOSE 8080
CMD ["./forum-server"]