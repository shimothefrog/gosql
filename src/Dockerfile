# FROM golang:1.23.0-alpine3.20
FROM golang:1.23.1-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# RUN go build -trimpath -ldflags="-s -w"
RUN go build -trimpath

# --- DEV STAGE ----
FROM builder AS dev
WORKDIR /src

COPY config/dev config/dev

RUN go install github.com/air-verse/air@latest

CMD ["air", "-c", "/src/config/dev/.air.toml"]

# --- PRD STAGE ----
FROM alpine:3.20.3
WORKDIR /root/
COPY --from=builder /src/sample .
CMD ["./sample"] 