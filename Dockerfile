# --- build stage ---
FROM golang:1.25-alpine AS builder
WORKDIR /app

# для CGO-сборки sqlite3
RUN apk add --no-cache build-base sqlite-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# включаем CGO
ENV CGO_ENABLED=1
RUN go build -o rest-api .

# --- runtime stage ---
FROM alpine:3.20
WORKDIR /app

# нужны рантайм-либы sqlite
RUN apk add --no-cache ca-certificates tzdata sqlite-libs

COPY --from=builder /app/rest-api .
EXPOSE 8080
ENV GIN_MODE=release
CMD ["./rest-api"]
